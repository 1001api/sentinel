package dto

import "time"

type CreateProjectInput struct {
	Name        string
	Description string
	UserID      string
	CreatedAt   time.Time
}
