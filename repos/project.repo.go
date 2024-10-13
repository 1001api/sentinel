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
	UpdateProject(ctx context.Context, name string, desc string, projectID string, userID string) error
	FindAll(ctx context.Context, userID string) ([]entities.Project, error)
	GetByID(ctx context.Context, projectID string, userID string) (*entities.Project, error)
	CheckWithinUserID(ctx context.Context, projectID string, userID string) (bool, error)
	CountProject(ctx context.Context, userID string) (int, error)
	DeleteProject(ctx context.Context, userID string, projectID string) error
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

func (r *ProjectRepositoryImpl) UpdateProject(ctx context.Context, name string, desc string, projectID string, userID string) error {
	SQL := "UPDATE projects SET name = $1, description = $2 WHERE id = $3 AND user_id = $4 AND deleted_at IS NULL"

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		log.Println("Failed preparing for transaction:", err)
		return err
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil {
				log.Println("Failed to rollback the tx:", err)
			}
		}
	}()

	_, err = tx.Exec(ctx, SQL, name, desc, projectID, userID)
	if err != nil {
		log.Println("Failed executing transaction:", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed committing tx:", err)
		return err
	}

	return nil
}

func (r *ProjectRepositoryImpl) FindAll(ctx context.Context, userID string) ([]entities.Project, error) {
	var projects []entities.Project

	SQL := "SELECT id, name, description, created_at FROM projects WHERE user_id = $1 AND deleted_at IS NULL"

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

func (r *ProjectRepositoryImpl) GetByID(ctx context.Context, projectID string, userID string) (*entities.Project, error) {
	var project entities.Project

	SQL := "SELECT id, name, description, created_at FROM projects WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL"
	row := r.DB.QueryRow(ctx, SQL, projectID, userID)
	if err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
	); err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepositoryImpl) CheckWithinUserID(ctx context.Context, projectID string, userID string) (bool, error) {
	var exist bool

	SQL := "SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL)"
	row := r.DB.QueryRow(ctx, SQL, projectID, userID)
	if err := row.Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

func (r *ProjectRepositoryImpl) CountProject(ctx context.Context, userID string) (int, error) {
	var total int

	SQL := "SELECT COUNT(*) FROM projects WHERE user_id = $1 AND deleted_at IS NULL"

	row := r.DB.QueryRow(ctx, SQL, userID)
	if err := row.Scan(
		&total,
	); err != nil {
		log.Println(err)
		return 0, err
	}

	return total, nil
}

func (r *ProjectRepositoryImpl) DeleteProject(ctx context.Context, userID string, projectID string) error {
	SQL := `
		UPDATE projects SET deleted_at = NOW() WHERE user_id = $1 AND id = $2 AND deleted_at IS NULL
	`

	DELETE_EVENT_SQL := "DELETE FROM events WHERE user_id = $1 AND project_id = $2"

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

	_, err = tx.Exec(ctx, DELETE_EVENT_SQL, userID, projectID)
	if err != nil {
		log.Println("Failed to delete associated events:", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit delete project tx:", err)
		return err
	}

	return nil
}
