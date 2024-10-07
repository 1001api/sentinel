package repositories

import (
	"context"
	"log"

	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProjectRepository interface {
	CreateProject(ctx context.Context, input *dto.CreateProjectInput) (*entities.Project, error)
	FindAll(ctx context.Context, userID string) ([]entities.Project, error)
	DeleteProject(ctx context.Context, userID string, projectID int) error
}

type ProjectRepositoryImpl struct {
	DB *pgxpool.Pool
}

func (r *ProjectRepositoryImpl) CreateProject(ctx context.Context, input *dto.CreateProjectInput) (*entities.Project, error) {
	var key entities.Project

	SQL := "INSERT INTO projects(name, description, user_id, created_at) VALUES ($1, $2, $3, $4) RETURNING name, description, created_at"

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Println("Failed preparing for transaction:", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Println("Failed to rollback the tx:", err)
			}
		}
	}()

	row := tx.QueryRow(ctx, SQL, input.Name, input.Description, input.UserID, input.CreatedAt)
	if err := row.Scan(
		&key.Name,
		&key.Description,
		&key.CreatedAt,
	); err != nil {
		log.Println("Failed creating project:", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed committing tx:", err)
		return nil, err
	}

	return &key, nil
}

func (r *ProjectRepositoryImpl) FindAll(ctx context.Context, userID string) ([]entities.Project, error) {
	var projects []entities.Project

	SQL := "SELECT id, name, description, created_at FROM projects WHERE user_id = $1"

	rows, err := r.DB.Query(ctx, SQL, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var project entities.Project

		if err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

func (r *ProjectRepositoryImpl) DeleteProject(ctx context.Context, userID string, projectID int) error {
	SQL := `
		DELETE FROM projects WHERE user_id = $1 AND id = $2
	`

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, SQL, userID, projectID)
	if err != nil {
		log.Println("Failed to delete project:", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit delete project tx:", err)
		return err
	}

	return nil
}
