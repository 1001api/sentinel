package services

import (
	"context"

	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/entities"
	repositories "github.com/hubkudev/sentinel/repos"
)

type UserService interface {
	FindByEmail(email string) (*entities.User, error)
	FindByID(userID string) (*entities.User, error)
	CreateUser(payload *dto.GooglePayload) (*entities.User, error)
}

type UserServiceImpl struct {
	UtilService UtilService
	UserRepo    repositories.UserRepository
}

func (s *UserServiceImpl) FindByEmail(email string) (*entities.User, error) {
	result, err := s.UserRepo.FindByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UserServiceImpl) FindByID(userID string) (*entities.User, error) {
	result, err := s.UserRepo.FindByID(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UserServiceImpl) CreateUser(payload *dto.GooglePayload) (*entities.User, error) {
	input := &dto.CreateUserInput{
		Fullname:      payload.GivenName,
		Email:         payload.Email,
		OAuthID:       payload.SUB,
		OAuthProvider: "google",
		ProfileURL:    payload.Picture,
	}

	result, err := s.UserRepo.CreateUser(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
