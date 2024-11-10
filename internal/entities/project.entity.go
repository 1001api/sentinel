package entities

import "time"

type Project struct {
	ID          string
	Name        string
	Description string
	UserID      string
	CreatedAt   time.Time
}
