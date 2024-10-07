package entities

import (
	"database/sql"
	"time"
)

type User struct {
	ID            string
	Fullname      string
	Email         string
	OAuthProvider string
	OAuthID       sql.NullString
	ProfileURL    sql.NullString
	CreatedAt     time.Time
	UpdatedAt     sql.NullTime
	DeletedAt     sql.NullTime
}
