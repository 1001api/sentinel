package services

import (
	"context"
	"time"

	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/entities"
	repositories "github.com/hubkudev/sentinel/repos"
)

type APIKeyService interface {
	CreateAPIKey(name string, userID string) (*entities.APIKey, error)
	GetAllKeys(userID string) ([]entities.APIKey, error)
}

type APIKeyServiceImpl struct {
	UtilService UtilService
	APIKeyRepo  repositories.APIKeyRepository
}

func (s *APIKeyServiceImpl) CreateAPIKey(name string, userID string) (*entities.APIKey, error) {
	input := dto.CreateAPIKeyInput{
		Name:      name,
		Token:     s.UtilService.GenerateRandomID(64),
		UserID:    userID,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().AddDate(0, 3, 0), // 3 months from now
	}

	key, err := s.APIKeyRepo.CreateKey(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s *APIKeyServiceImpl) GetAllKeys(userID string) ([]entities.APIKey, error) {
	key, err := s.APIKeyRepo.FindAll(context.Background(), userID)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (s *APIKeyServiceImpl) DeleteKey(userID string, keyID int) error {
	err := s.APIKeyRepo.DeleteKey(context.Background(), userID, keyID)
	if err != nil {
		return err
	}
	return nil
}
