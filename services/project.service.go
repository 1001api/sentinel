package services

import (
	"context"
	"errors"
	"time"

	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/entities"
	repositories "github.com/hubkudev/sentinel/repos"
)

type ProjectService interface {
	CreateProject(name string, desc string, userID string) (*entities.Project, error)
	UpdateProject(name string, desc string, projectID int, userID string) error
	GetAllProjects(userID string) ([]entities.Project, error)
	GetProjectCount(userID string) (int, error)
	DeleteProject(userID string, keyID int) error
}

type ProjectServiceImpl struct {
	ProjectRepo repositories.ProjectRepository
}

func (s *ProjectServiceImpl) CreateProject(name string, desc string, userID string) (*entities.Project, error) {
	// check how many projects already this user has
	count, err := s.ProjectRepo.CountProject(context.Background(), userID)
	if err != nil {
		return nil, err
	}

	// if project is more or equal to 5
	// reject creation.
	if count >= 5 {
		return nil, errors.New("Total project already at max") // reject with error
	}

	input := dto.CreateProjectInput{
		Name:        name,
		Description: desc,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	key, err := s.ProjectRepo.CreateProject(context.Background(), &input)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s *ProjectServiceImpl) UpdateProject(name string, desc string, projectID int, userID string) error {
	return s.ProjectRepo.UpdateProject(context.Background(), name, desc, projectID, userID)
}

func (s *ProjectServiceImpl) GetAllProjects(userID string) ([]entities.Project, error) {
	return s.ProjectRepo.FindAll(context.Background(), userID)
}

func (s *ProjectServiceImpl) GetProjectCount(userID string) (int, error) {
	return s.ProjectRepo.CountProject(context.Background(), userID)
}

func (s *ProjectServiceImpl) DeleteProject(userID string, keyID int) error {
	return s.ProjectRepo.DeleteProject(context.Background(), userID, keyID)
}
