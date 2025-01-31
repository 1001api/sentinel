package repositories

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepo interface {
	Create(ctx context.Context, input *gen.CreateProjectParams) (gen.CreateProjectRow, error)
	Update(ctx context.Context, input *gen.UpdateProjectParams) error
	FindByID(ctx context.Context, input *gen.FindProjectByIDParams) (gen.FindProjectByIDRow, error)
	FindAll(ctx context.Context, userID uuid.UUID) ([]gen.FindAllProjectsRow, error)
	Count(ctx context.Context, userID uuid.UUID) (int64, error)
	CountSize(ctx context.Context, input *gen.CountProjectSizeParams) (int64, error)
	Delete(ctx context.Context, input *gen.DeleteProjectParams) error
	LastDataReceived(ctx context.Context, input *gen.LastProjectDataReceivedParams) (time.Time, error)
	CheckProjectAggrEligibility(ctx context.Context, projectID uuid.UUID) (int64, error)
	CreateProjectAggr(ctx context.Context, input *gen.CreateProjectAggrParams) error
	FindProjectAggr(ctx context.Context, input *gen.FindProjectAggrParams) ([]gen.ProjectAggregation, error)
}

type ProjectRepoImpl struct {
	Repo *gen.Queries
	DB   *pgxpool.Pool
}

func InitProjectRepo(repo *gen.Queries, db *pgxpool.Pool) ProjectRepoImpl {
	return ProjectRepoImpl{
		Repo: repo,
		DB:   db,
	}
}

func (r *ProjectRepoImpl) Create(ctx context.Context, input *gen.CreateProjectParams) (gen.CreateProjectRow, error) {
	return r.Repo.CreateProject(ctx, *input)
}

func (r *ProjectRepoImpl) Update(ctx context.Context, input *gen.UpdateProjectParams) error {
	return r.Repo.UpdateProject(ctx, *input)
}

func (r *ProjectRepoImpl) FindByID(ctx context.Context, input *gen.FindProjectByIDParams) (gen.FindProjectByIDRow, error) {
	return r.Repo.FindProjectByID(ctx, *input)
}

func (r *ProjectRepoImpl) FindAll(ctx context.Context, userID uuid.UUID) ([]gen.FindAllProjectsRow, error) {
	return r.Repo.FindAllProjects(ctx, userID)
}

func (r *ProjectRepoImpl) Count(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.Repo.CountProject(ctx, userID)
}

func (r *ProjectRepoImpl) CountSize(ctx context.Context, input *gen.CountProjectSizeParams) (int64, error) {
	return r.Repo.CountProjectSize(ctx, *input)
}

func (r *ProjectRepoImpl) Delete(ctx context.Context, input *gen.DeleteProjectParams) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Println("error starting the transaction", err)
		return err
	}
	defer tx.Rollback(ctx)

	qtx := r.Repo.WithTx(tx)

	// delete project
	if err = qtx.DeleteProject(ctx, *input); err != nil {
		return err
	}

	// delete event related to project
	if err = qtx.DeleteEventByProjectID(ctx, gen.DeleteEventByProjectIDParams{
		UserID:    input.UserID,
		ProjectID: input.ID,
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

func (r *ProjectRepoImpl) LastDataReceived(ctx context.Context, input *gen.LastProjectDataReceivedParams) (time.Time, error) {
	return r.Repo.LastProjectDataReceived(ctx, *input)
}

func (r *ProjectRepoImpl) CheckProjectAggrEligibility(ctx context.Context, projectID uuid.UUID) (int64, error) {
	return r.Repo.CheckProjectAggrEligibility(ctx, projectID)
}

func (r *ProjectRepoImpl) CreateProjectAggr(ctx context.Context, input *gen.CreateProjectAggrParams) error {
	return r.Repo.CreateProjectAggr(ctx, *input)
}

func (r *ProjectRepoImpl) FindProjectAggr(ctx context.Context, input *gen.FindProjectAggrParams) ([]gen.ProjectAggregation, error) {
	return r.Repo.FindProjectAggr(ctx, *input)
}
