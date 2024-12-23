package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/repositories"
	"github.com/jackc/pgx/v5/pgtype"
)

type KeyService interface {
	CreateAPIKey(ctx context.Context, name string, userID uuid.UUID) (*gen.CreateAPIKeyRow, error)
	GetAllKeys(ctx context.Context, userID uuid.UUID) ([]gen.FindAllAPIKeysRow, error)
	DeleteKey(ctx context.Context, userID uuid.UUID, keyID int) error
}

type KeyServiceImpl struct {
	UtilService UtilService
	Repo        repositories.KeyRepo
}

func InitKeyService(
	utilService UtilService,
	repo repositories.KeyRepo,
) KeyServiceImpl {
	return KeyServiceImpl{
		UtilService: utilService,
		Repo:        repo,
	}
}

func (s *KeyServiceImpl) CreateAPIKey(ctx context.Context, name string, userID uuid.UUID) (*gen.CreateAPIKeyRow, error) {
	input := gen.CreateAPIKeyParams{
		Name:   name,
		Token:  fmt.Sprintf("snt_%s", s.UtilService.GenerateRandomID(48)),
		UserID: userID,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
		ExpiredAt: pgtype.Timestamptz{
			Time:  time.Now().AddDate(0, 6, 0), // 6 months from now
			Valid: true,
		},
	}

	key, err := s.Repo.CreateAPIKey(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (s *KeyServiceImpl) GetAllKeys(ctx context.Context, userID uuid.UUID) ([]gen.FindAllAPIKeysRow, error) {
	key, err := s.Repo.GetAllPrivateKeys(ctx, userID)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *KeyServiceImpl) DeleteKey(ctx context.Context, userID uuid.UUID, keyID int) error {
	err := s.Repo.DeletePrivateKey(ctx, &gen.DeleteAPIKeyParams{
		UserID: userID,
		ID:     int32(keyID),
	})
	if err != nil {
		return err
	}
	return nil
}
