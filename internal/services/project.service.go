package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/repositories"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectService interface {
	CreateProject(ctx context.Context, name string, desc string, userID uuid.UUID) (*gen.CreateProjectRow, error)
	UpdateProject(ctx context.Context, name string, desc string, projectID uuid.UUID, userID uuid.UUID) error
	GetProjectByID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*gen.FindProjectByIDRow, error)
	GetAllProjects(ctx context.Context, userID uuid.UUID) ([]gen.FindAllProjectsRow, error)
	GetProjectCount(ctx context.Context, userID uuid.UUID) (int64, error)
	DeleteProject(ctx context.Context, userID uuid.UUID, projectID uuid.UUID) error
	CountProjectSize(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (int64, error)
	LastProjectDataReceived(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*time.Time, error)
}

type ProjectServiceImpl struct {
	Repo repositories.ProjectRepo
}

func InitProjectService(repo repositories.ProjectRepo) ProjectServiceImpl {
	return ProjectServiceImpl{
		Repo: repo,
	}
}

func (s *ProjectServiceImpl) CreateProject(ctx context.Context, name string, desc string, userID uuid.UUID) (*gen.CreateProjectRow, error) {
	input := gen.CreateProjectParams{
		Name: name,
		Description: pgtype.Text{
			String: desc,
			Valid:  desc != "",
		},
		UserID: userID,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	key, err := s.Repo.Create(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (s *ProjectServiceImpl) UpdateProject(ctx context.Context, name string, desc string, projectID uuid.UUID, userID uuid.UUID) error {
	input := gen.UpdateProjectParams{
		Name: name,
		Description: pgtype.Text{
			String: desc,
			Valid:  desc != "",
		},
		ID:     projectID,
		UserID: userID,
	}
	return s.Repo.Update(ctx, &input)
}

func (s *ProjectServiceImpl) GetProjectByID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*gen.FindProjectByIDRow, error) {
	input := gen.FindProjectByIDParams{
		ID:     projectID,
		UserID: userID,
	}

	row, err := s.Repo.FindByID(ctx, &input)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (s *ProjectServiceImpl) GetAllProjects(ctx context.Context, userID uuid.UUID) ([]gen.FindAllProjectsRow, error) {
	return s.Repo.FindAll(ctx, userID)
}

func (s *ProjectServiceImpl) GetProjectCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.Repo.Count(ctx, userID)
}

func (s *ProjectServiceImpl) DeleteProject(ctx context.Context, userID uuid.UUID, projectID uuid.UUID) error {
	input := gen.DeleteProjectParams{
		UserID: userID,
		ID:     projectID,
	}
	return s.Repo.Delete(ctx, &input)
}

func (s *ProjectServiceImpl) CountProjectSize(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (int64, error) {
	input := gen.CountProjectSizeParams{
		UserID:    userID,
		ProjectID: projectID,
	}

	size, err := s.Repo.CountSize(ctx, &input)
	if err != nil {
		return -1, err
	}
	return size, err
}

func (s *ProjectServiceImpl) LastProjectDataReceived(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*time.Time, error) {
	input := gen.LastProjectDataReceivedParams{
		UserID:    userID,
		ProjectID: projectID,
	}

	lastTime, err := s.Repo.LastDataReceived(ctx, &input)
	if err != nil {
		return nil, err
	}
	return &lastTime, err
}
