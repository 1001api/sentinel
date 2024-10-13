package repositories

import (
	"context"
	"log"
	"time"

	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, input *dto.CreateEventInput, userID string) error
	GetLiveEvents(ctx context.Context, userID string) ([]entities.Event, error)
	GetEventSummary(ctx context.Context, projectID string, userID string) (*entities.EventSummary, error)
}

type EventRepoImpl struct {
	DB *pgxpool.Pool
}

func (r *EventRepoImpl) CreateEvent(ctx context.Context, input *dto.CreateEventInput, userID string) error {
	var ipAddr *string

	// if ip address is not specified
	if input.IPAddr != "" {
		ipAddr = &input.IPAddr
	}

	SQL := `
		INSERT INTO events (
			event_type,
			event_label,
			page_url,
			element_path,
			element_type,
			ip_addr,
			user_agent,
			browser_name,
			country,
			region,
			city,
			session_id,
			device_type,
			time_on_page,
			screen_resolution,
			fired_at,
			received_at,
			user_id,
			project_id
		) VALUES (
			$1,  -- event_type
			$2,  -- event_label
			$3,  -- page_url
			$4,  -- element_path
			$5,  -- element_type
			$6,  -- ip_addr
			$7,  -- user_agent
			$8,  -- browser_name
			$9,  -- country
			$10, -- region
			$11, -- city
			$12, -- session_id
			$13, -- device_type
			$14, -- time_on_page
			$15, -- screen_resolution
			$16, -- fired_at
			$17, -- received_at
			$18, -- user_id
			$19  -- project_id
		);
	`

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Println("Failed preparing for transaction:", err)
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Println("Failed to rollback the tx:", err)
			}
		}
	}()

	if _, err := tx.Exec(ctx,
		SQL,
		input.EventType,
		input.EventLabel,
		input.PageURL,
		input.ElementPath,
		input.ElementType,
		ipAddr,
		input.UserAgent,
		input.BrowserName,
		input.Country,
		input.Region,
		input.City,
		input.SessionID,
		input.DeviceType,
		input.TimeOnPage,
		input.ScreenResolution,
		input.FiredAt,
		time.Now(),
		userID,
		input.ProjectID,
	); err != nil {
		log.Println("Failed creating event:", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed committing tx:", err)
		return err
	}

	return nil
}

func (r *EventRepoImpl) GetLiveEvents(ctx context.Context, userID string) ([]entities.Event, error) {
	var events []entities.Event

	SQL := `
		SELECT
			p.name,
			e.event_type,
			e.event_label,
			e.page_url,
			e.element_path,
			e.country,
			e.fired_at,
			e.received_at
		FROM events AS e
		JOIN projects AS p ON e.project_id = p.id
		WHERE e.user_id = $1 AND e.received_at >= NOW() - INTERVAL '1 hour'
		ORDER BY e.received_at DESC
		LIMIT 100;
	`

	rows, err := r.DB.Query(ctx, SQL, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e entities.Event

		if err = rows.Scan(
			&e.ProjectName,
			&e.EventType,
			&e.EventLabel,
			&e.PageURL,
			&e.ElementPath,
			&e.Country,
			&e.FiredAt,
			&e.ReceivedAt,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		events = append(events, e)
	}

	return events, nil
}

func (r *EventRepoImpl) GetEventSummary(ctx context.Context, projectID string, userID string) (*entities.EventSummary, error) {
	var summary entities.EventSummary

	SQL := `
		SELECT 
		COUNT(id) AS total_events,
		COUNT(DISTINCT ip_addr) AS total_unique_users,
		COUNT(DISTINCT event_type) AS total_event_type,
		COUNT(DISTINCT country) AS total_country_visited,
		(
			SELECT page_url FROM events
			WHERE user_id = $2 AND project_id = $1
			GROUP BY page_url ORDER BY COUNT(page_url) DESC LIMIT 1
		) AS most_visited_url,
		(
			SELECT country FROM events
			WHERE user_id = $2 AND project_id = $1
			GROUP BY country ORDER BY COUNT(country) DESC LIMIT 1
		) AS most_country_visited
		FROM events WHERE user_id = $2 AND project_id = $1;
	`

	row := r.DB.QueryRow(ctx, SQL, projectID, userID)
	if err := row.Scan(
		&summary.TotalEvents,
		&summary.TotalUniqueUsers,
		&summary.TotalEventType,
		&summary.TotalCountryVisited,
		&summary.MostVisitedURL,
		&summary.MostCountryVisited,
	); err != nil {
		log.Println(err)
		return nil, err
	}

	return &summary, nil
}
