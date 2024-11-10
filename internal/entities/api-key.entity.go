package entities

import "time"

type APIKey struct {
	ID        int
	Name      string
	Token     string
	UserID    string
	CreatedAt time.Time
	ExpiredAt time.Time
}
