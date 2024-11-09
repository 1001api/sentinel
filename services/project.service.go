package services

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectService interface {
	CreateProject(ctx context.Context, name string, desc string, userID string) (*gen.CreateProjectRow, error)
	UpdateProject(ctx context.Context, name string, desc string, projectID string, userID string) error
	GetProjectByID(ctx context.Context, projectID string, userID string) (*gen.FindProjectByIDRow, error)
	GetAllProjects(ctx context.Context, userID string) ([]gen.FindAllProjectsRow, error)
	GetProjectCount(ctx context.Context, userID string) (int64, error)
	DeleteProject(ctx context.Context, userID string, projectID string) error
	CountProjectSize(ctx context.Context, projectID string, userID string) (int64, error)
	LastProjectDataReceived(ctx context.Context, projectID string, userID string) (*time.Time, error)
}

type ProjectServiceImpl struct {
	Repo *gen.Queries
	DB   *pgxpool.Pool
}

func InitProjectService(repo *gen.Queries, db *pgxpool.Pool) ProjectServiceImpl {
	return ProjectServiceImpl{
		Repo: repo,
		DB:   db,
	}
}

func (s *ProjectServiceImpl) CreateProject(ctx context.Context, name string, desc string, userID string) (*gen.CreateProjectRow, error) {
	userUUID := uuid.MustParse(userID)

	input := gen.CreateProjectParams{
		Name: name,
		Description: pgtype.Text{
			String: desc,
			Valid:  desc != "",
		},
		UserID: userUUID,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	key, err := s.Repo.CreateProject(ctx, input)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (s *ProjectServiceImpl) UpdateProject(ctx context.Context, name string, desc string, projectID string, userID string) error {
	userUUID, projectUUID := uuid.MustParse(userID), uuid.MustParse(projectID)
	return s.Repo.UpdateProject(ctx, gen.UpdateProjectParams{
		Name: name,
		Description: pgtype.Text{
			String: desc,
			Valid:  desc != "",
		},
		ID:     projectUUID,
		UserID: userUUID,
	})
}

func (s *ProjectServiceImpl) GetProjectByID(ctx context.Context, projectID string, userID string) (*gen.FindProjectByIDRow, error) {
	userUUID, projectUUID := uuid.MustParse(userID), uuid.MustParse(projectID)
	row, err := s.Repo.FindProjectByID(ctx, gen.FindProjectByIDParams{
		ID:     projectUUID,
		UserID: userUUID,
	})
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (s *ProjectServiceImpl) GetAllProjects(ctx context.Context, userID string) ([]gen.FindAllProjectsRow, error) {
	userUUID := uuid.MustParse(userID)
	return s.Repo.FindAllProjects(ctx, userUUID)
}

func (s *ProjectServiceImpl) GetProjectCount(ctx context.Context, userID string) (int64, error) {
	userUUID := uuid.MustParse(userID)
	return s.Repo.CountProject(ctx, userUUID)
}

func (s *ProjectServiceImpl) DeleteProject(ctx context.Context, userID string, projectID string) error {
	userUUID, projectUUID := uuid.MustParse(userID), uuid.MustParse(projectID)

	tx, err := s.DB.Begin(ctx)
	if err != nil {
		log.Println("error starting the transaction", err)
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.Repo.WithTx(tx)

	// delete project
	if err = qtx.DeleteProject(ctx, gen.DeleteProjectParams{
		UserID: userUUID,
		ID:     projectUUID,
	}); err != nil {
		return err
	}

	// delete event related to project
	if err = qtx.DeleteEventByProjectID(ctx, gen.DeleteEventByProjectIDParams{
		UserID:    userUUID,
		ProjectID: projectUUID,
	}); err != nil {
		return err
	}

	// commit if everything alright
	if err = tx.Commit(ctx); err != nil {
		log.Println("error commiting the transaction", err)
		return err
	}

	return nil
}

func (s *ProjectServiceImpl) CountProjectSize(ctx context.Context, projectID string, userID string) (int64, error) {
	userUUID, projectUUID := uuid.MustParse(userID), uuid.MustParse(projectID)

	size, err := s.Repo.CountProjectSize(ctx, gen.CountProjectSizeParams{
		UserID:    userUUID,
		ProjectID: projectUUID,
	})
	if err != nil {
		return -1, err
	}

	return size, err
}

func (s *ProjectServiceImpl) LastProjectDataReceived(ctx context.Context, projectID string, userID string) (*time.Time, error) {
	userUUID, projectUUID := uuid.MustParse(userID), uuid.MustParse(projectID)

	lastTime, err := s.Repo.LastProjectDataReceived(ctx, gen.LastProjectDataReceivedParams{
		UserID:    userUUID,
		ProjectID: projectUUID,
	})
	if err != nil {
		return nil, err
	}

	return &lastTime, err
}
