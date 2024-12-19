package entities

import "time"

type Project struct {
	ID          string
	Name        string
	Url         string
	Description string
	UserID      string
	CreatedAt   time.Time
}
