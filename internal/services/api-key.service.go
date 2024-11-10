package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/repositories"
	"github.com/jackc/pgx/v5/pgtype"
)

type APIKeyService interface {
	CreateAPIKey(name string, userID uuid.UUID) (*gen.CreateAPIKeyRow, error)
	GetAllKeys(userID uuid.UUID) ([]gen.FindAllAPIKeysRow, error)
	DeleteKey(userID uuid.UUID, keyID int) error
}

type APIKeyServiceImpl struct {
	UtilService UtilService
	Repo        repositories.APIKeyRepoImpl
}

func (s *APIKeyServiceImpl) CreateAPIKey(name string, userID uuid.UUID) (*gen.CreateAPIKeyRow, error) {
	input := gen.CreateAPIKeyParams{
		Name:   name,
		Token:  s.UtilService.GenerateRandomID(64),
		UserID: userID,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ExpiredAt: pgtype.Timestamptz{
			Time:  time.Now().AddDate(0, 3, 0), // 3 months from now
			Valid: true,
		},
	}

	key, err := s.Repo.CreateAPIKey(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (s *APIKeyServiceImpl) GetAllKeys(userID uuid.UUID) ([]gen.FindAllAPIKeysRow, error) {
	key, err := s.Repo.GetAllPrivateKeys(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *APIKeyServiceImpl) DeleteKey(userID uuid.UUID, keyID int) error {
	err := s.Repo.DeletePrivateKey(context.Background(), &gen.DeleteAPIKeyParams{
		UserID: userID,
		ID:     int32(keyID),
	})
	if err != nil {
		return err
	}
	return nil
}
