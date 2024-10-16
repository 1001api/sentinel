package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type APIKeyService interface {
	CreateAPIKey(name string, userID string) (*gen.CreateAPIKeyRow, error)
	GetAllKeys(userID string) ([]gen.FindAllAPIKeysRow, error)
	DeleteKey(userID string, keyID int) error
}

type APIKeyServiceImpl struct {
	UtilService UtilService
	Repo        *gen.Queries
}

func (s *APIKeyServiceImpl) CreateAPIKey(name string, userID string) (*gen.CreateAPIKeyRow, error) {
	userUUID := uuid.MustParse(userID)

	input := gen.CreateAPIKeyParams{
		Name:   name,
		Token:  s.UtilService.GenerateRandomID(64),
		UserID: userUUID,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ExpiredAt: pgtype.Timestamptz{
			Time:  time.Now().AddDate(0, 3, 0), // 3 months from now
			Valid: true,
		},
	}

	key, err := s.Repo.CreateAPIKey(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (s *APIKeyServiceImpl) GetAllKeys(userID string) ([]gen.FindAllAPIKeysRow, error) {
	userUUID := uuid.MustParse(userID)
	key, err := s.Repo.FindAllAPIKeys(context.Background(), userUUID)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *APIKeyServiceImpl) DeleteKey(userID string, keyID int) error {
	userUUID := uuid.MustParse(userID)
	err := s.Repo.DeleteAPIKey(context.Background(), gen.DeleteAPIKeyParams{
		UserID: userUUID,
		ID:     int32(keyID),
	})
	if err != nil {
		return err
	}
	return nil
}
