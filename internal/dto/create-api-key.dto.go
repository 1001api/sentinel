package dto

import "time"

type CreateAPIKeyInput struct {
	Name      string
	Token     string
	UserID    string
	CreatedAt time.Time
	ExpiredAt time.Time
}
