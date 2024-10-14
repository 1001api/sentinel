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
	GetLiveEventDetail(ctx context.Context, projectID string, userID string) ([]entities.Event, error)
	GetEventSummary(ctx context.Context, projectID string, userID string) (*entities.EventSummary, error)
	GetEventDetailSummary(ctx context.Context, projectID string, userID string) (*entities.EventDetail, error)
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

func (r *EventRepoImpl) GetLiveEventDetail(ctx context.Context, projectID string, userID string) ([]entities.Event, error) {
	var events []entities.Event

	SQL := `
		SELECT 
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
			device_type,
			time_on_page,
			screen_resolution,
			fired_at,
			received_at
		FROM events
		WHERE user_id = $2 AND project_id = $1 AND received_at >= NOW() - INTERVAL '1 hour'
		ORDER BY received_at DESC
		LIMIT 50;
	`

	rows, err := r.DB.Query(ctx, SQL, projectID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e entities.Event

		if err = rows.Scan(
			&e.EventType,
			&e.EventLabel,
			&e.PageURL,
			&e.ElementPath,
			&e.ElementType,
			&e.IPAddr,
			&e.UserAgent,
			&e.BrowserName,
			&e.Country,
			&e.Region,
			&e.City,
			&e.DeviceType,
			&e.TimeOnPage,
			&e.ScreenResolution,
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

func (r *EventRepoImpl) GetEventDetailSummary(ctx context.Context, projectID string, userID string) (*entities.EventDetail, error) {
	var summary entities.EventDetail

	TOTAL_SQL := `
		SELECT 
		COUNT(id) AS total_events,
		COUNT(DISTINCT event_type) AS total_event_type,
		COUNT(DISTINCT ip_addr) AS total_unique_users,
		COUNT(DISTINCT country) AS total_country_visited,
		COUNT(DISTINCT page_url) AS total_page_url
		FROM events WHERE user_id = $2 AND project_id = $1;
	`
	row := r.DB.QueryRow(ctx, TOTAL_SQL, projectID, userID)
	if err := row.Scan(
		&summary.TotalEvents,
		&summary.TotalEventType,
		&summary.TotalUniqueUsers,
		&summary.TotalCountryVisited,
		&summary.TotalPageURL,
	); err != nil {
		log.Println(err)
		return nil, err
	}

	URL_SQL := `
		SELECT page_url AS name, COUNT(page_url) AS total
		FROM events WHERE user_id = $2 AND project_id = $1
		GROUP BY page_url ORDER BY COUNT(page_url) DESC LIMIT 5
	`
	rows, err := r.DB.Query(ctx, URL_SQL, projectID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e entities.EventTextTotal

		if err = rows.Scan(
			&e.Name,
			&e.Total,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		summary.MostVisitedURLs = append(summary.MostVisitedURLs, e)
	}

	COUNTRIES_SQL := `
		SELECT country AS name, COUNT(country) AS total
		FROM events 
		WHERE country IS NOT NULL AND country <> ''
		AND user_id = $2 AND project_id = $1
		GROUP BY country ORDER BY COUNT(country) DESC LIMIT 5
	`
	rows, err = r.DB.Query(ctx, COUNTRIES_SQL, projectID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var e entities.EventTextTotal

		if err = rows.Scan(
			&e.Name,
			&e.Total,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		summary.MostCountryVisited = append(summary.MostCountryVisited, e)
	}

	CITIES_SQL := `
		SELECT city AS name, COUNT(city) AS total
		FROM events 
		WHERE city IS NOT NULL AND city <> '' 
		AND user_id = $2 AND project_id = $1
		GROUP BY city ORDER BY COUNT(city) DESC LIMIT 5
	`
	rows, err = r.DB.Query(ctx, CITIES_SQL, projectID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var e entities.EventTextTotal

		if err = rows.Scan(
			&e.Name,
			&e.Total,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		summary.MostCitiesVisited = append(summary.MostCitiesVisited, e)
	}

	ELEMENTS_SQL := `
		SELECT element_path AS name, COUNT(element_path) AS total
		FROM events 
		WHERE element_path IS NOT NULL AND element_path <> '' 
		AND user_id = $2 AND project_id = $1
		GROUP BY element_path ORDER BY COUNT(element_path) DESC LIMIT 5
	`
	rows, err = r.DB.Query(ctx, ELEMENTS_SQL, projectID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var e entities.EventTextTotal

		if err = rows.Scan(
			&e.Name,
			&e.Total,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		summary.MostElementsFired = append(summary.MostElementsFired, e)
	}

	LAST_USERS_SQL := `
		SELECT ip_addr AS ip, received_at AS timestamp
		FROM events
		WHERE user_id = $2 AND project_id = $1
		GROUP BY ip_addr, received_at ORDER BY received_at DESC LIMIT 5
	`
	rows, err = r.DB.Query(ctx, LAST_USERS_SQL, projectID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var v entities.EventLastUser

		if err = rows.Scan(
			&v.IP,
			&v.Timestamp,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		summary.LastVisitedUsers = append(summary.LastVisitedUsers, v)
	}

	BROWSERS_SQL := `
		SELECT browser_name AS name, COUNT(browser_name) AS total
		FROM events
		WHERE browser_name IS NOT NULL AND browser_name <> ''
		AND user_id = $2 AND project_id = $1
		GROUP BY browser_name ORDER BY COUNT(browser_name) DESC LIMIT 5
	`
	rows, err = r.DB.Query(ctx, BROWSERS_SQL, projectID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var e entities.EventTextTotal

		if err = rows.Scan(
			&e.Name,
			&e.Total,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		summary.MostUsedBrowsers = append(summary.MostUsedBrowsers, e)
	}

	EVENT_TYPE_SQL := `
		SELECT event_type AS name, COUNT(event_type) AS total
		FROM events
		WHERE event_type IS NOT NULL AND event_type <> ''
		AND user_id = $2 AND project_id = $1
		GROUP BY event_type ORDER BY COUNT(event_type) DESC LIMIT 5
	`
	rows, err = r.DB.Query(ctx, EVENT_TYPE_SQL, projectID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var e entities.EventTextTotal

		if err = rows.Scan(
			&e.Name,
			&e.Total,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		summary.MostFiredEventType = append(summary.MostFiredEventType, e)
	}

	EVENT_LABEL_SQL := `
		SELECT event_label AS name, COUNT(event_label) AS total
		FROM events
		WHERE event_label IS NOT NULL AND event_label <> ''
		AND user_id = $2 AND project_id = $1
		GROUP BY event_label ORDER BY COUNT(event_label) DESC LIMIT 5
	`
	rows, err = r.DB.Query(ctx, EVENT_LABEL_SQL, projectID, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for rows.Next() {
		var e entities.EventTextTotal

		if err = rows.Scan(
			&e.Name,
			&e.Total,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		summary.MostFiredEventLabel = append(summary.MostFiredEventLabel, e)
	}

	return &summary, nil
}
