package services

import (
	"context"

	"github.com/google/uuid"
	gen "github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/repositories"
)

type UserService interface {
	FindByEmail(email string) (*gen.FindUserByEmailRow, error)
	FindByEmailWithHash(email string) (*gen.FindUserByEmailWithHashRow, error)
	FindByID(userID uuid.UUID) (*gen.FindUserByIDRow, error)
	FindByPublicKey(key string) (*gen.FindUserByPublicKeyRow, error)
	FindByPrivateKey(key string) (*gen.FindUserByPrivateKeyRow, error)
	CheckAdminExist() (bool, error)
	GetPublicKey(userID uuid.UUID) (string, error)
	CreateUser(payload *gen.CreateUserParams) (*gen.CreateUserRow, error)
}

type UserServiceImpl struct {
	UtilService UtilService
	Repo        repositories.UserRepo
}

func InitUserService(utilService UtilService, repo repositories.UserRepo) UserServiceImpl {
	return UserServiceImpl{
		UtilService: utilService,
		Repo:        repo,
	}
}

func (s *UserServiceImpl) FindByEmail(email string) (*gen.FindUserByEmailRow, error) {
	result, err := s.Repo.FindUserByEmail(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *UserServiceImpl) FindByEmailWithHash(email string) (*gen.FindUserByEmailWithHashRow, error) {
	result, err := s.Repo.FindUserByEmailWithHash(context.Background(), email)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *UserServiceImpl) FindByID(userID uuid.UUID) (*gen.FindUserByIDRow, error) {
	result, err := s.Repo.FindUserByID(context.Background(), userID)
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

func (s *UserServiceImpl) FindByPrivateKey(key string) (*gen.FindUserByPrivateKeyRow, error) {
	result, err := s.Repo.FindUserByPrivateKey(context.Background(), key)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *UserServiceImpl) GetPublicKey(userID uuid.UUID) (string, error) {
	result, err := s.Repo.FindUserPublicKey(context.Background(), userID)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (s *UserServiceImpl) CheckAdminExist() (bool, error) {
	return s.Repo.CheckAdminExist(context.Background())
}

func (s *UserServiceImpl) CreateUser(payload *gen.CreateUserParams) (*gen.CreateUserRow, error) {
	// generate random 48 long for public key
	key := s.UtilService.GenerateRandomID(48)
	payload.PublicKey = key

	result, err := s.Repo.CreateUser(context.Background(), payload)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
