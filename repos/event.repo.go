package repositories

import (
	"context"
	"log"
	"time"

	"github.com/hubkudev/sentinel/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, input *dto.CreateEventInput, userID string) error
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
