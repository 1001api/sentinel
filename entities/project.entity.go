package entities

import "time"

type Project struct {
	ID          int
	Name        string
	Description string
	UserID      string
	CreatedAt   time.Time
}
