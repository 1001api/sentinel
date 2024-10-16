package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/dto"
	gen "github.com/hubkudev/sentinel/gen"
	"github.com/jackc/pgx/v5/pgtype"
)

type UserService interface {
	FindByEmail(email string) (*gen.FindUserByEmailRow, error)
	FindByID(userID string) (*gen.FindUserByIDRow, error)
	FindByPublicKey(userID string) (*gen.FindUserByPublicKeyRow, error)
	GetPublicKey(userID string) (string, error)
	CreateUser(payload *dto.GooglePayload) (*gen.CreateUserRow, error)
}

type UserServiceImpl struct {
	UtilService UtilService
	Repo        *gen.Queries
}

func (s *UserServiceImpl) FindByEmail(email string) (*gen.FindUserByEmailRow, error) {
	result, err := s.Repo.FindUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *UserServiceImpl) FindByID(userID string) (*gen.FindUserByIDRow, error) {
	id := uuid.MustParse(userID)
	result, err := s.Repo.FindUserByID(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *UserServiceImpl) FindByPublicKey(key string) (*gen.FindUserByPublicKeyRow, error) {
	result, err := s.Repo.FindUserByPublicKey(context.Background(), key)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *UserServiceImpl) GetPublicKey(userID string) (string, error) {
	id := uuid.MustParse(userID)
	result, err := s.Repo.FindUserPublicKey(context.Background(), id)
	if err != nil {
		return "", err
	}

	return result, nil
}

func (s *UserServiceImpl) CreateUser(payload *dto.GooglePayload) (*gen.CreateUserRow, error) {
	// generate random 48 long for public key
	key := s.UtilService.GenerateRandomID(48)

	input := gen.CreateUserParams{
		Fullname: payload.GivenName,
		Email:    payload.Email,
		OauthID: pgtype.Text{
			String: payload.SUB,
		},
		OauthProvider: "google",
		ProfileUrl: pgtype.Text{
			String: payload.Picture,
		},
		PublicKey: key,
	}

	result, err := s.Repo.CreateUser(context.Background(), input)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
