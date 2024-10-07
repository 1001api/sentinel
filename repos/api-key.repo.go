package repositories

import (
	"context"
	"log"

	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIKeyRepository interface {
	CreateKey(ctx context.Context, input *dto.CreateAPIKeyInput) (*entities.APIKey, error)
	FindAll(ctx context.Context, userID string) ([]entities.APIKey, error)
	DeleteKey(ctx context.Context, userID string, keyID int) error
}

type APIKeyRepoImpl struct {
	DB *pgxpool.Pool
}

func (r *APIKeyRepoImpl) CreateKey(ctx context.Context, input *dto.CreateAPIKeyInput) (*entities.APIKey, error) {
	var key entities.APIKey

	SQL := "INSERT INTO api_keys(name, token, user_id, created_at, expired_at) VALUES ($1, $2, $3, $4, $5) RETURNING name, token, created_at, expired_at"

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

	row := tx.QueryRow(ctx, SQL, input.Name, input.Token, input.UserID, input.CreatedAt, input.ExpiredAt)
	if err := row.Scan(
		&key.Name,
		&key.Token,
		&key.CreatedAt,
		&key.ExpiredAt,
	); err != nil {
		log.Println("Failed creating api key:", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed committing tx:", err)
		return nil, err
	}

	return &key, nil
}

func (r *APIKeyRepoImpl) FindAll(ctx context.Context, userID string) ([]entities.APIKey, error) {
	var keys []entities.APIKey

	SQL := "SELECT id, name, token, created_at, expired_at FROM api_keys WHERE user_id = $1"

	rows, err := r.DB.Query(ctx, SQL, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var key entities.APIKey

		if err = rows.Scan(
			&key.ID,
			&key.Name,
			&key.Token,
			&key.CreatedAt,
			&key.ExpiredAt,
		); err != nil {
			log.Println(err)
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func (r *APIKeyRepoImpl) DeleteKey(ctx context.Context, userID string, keyID int) error {
	SQL := `
		DELETE FROM api_keys WHERE user_id = $1 AND id = $2
	`

	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, SQL, userID, keyID)
	if err != nil {
		log.Println("Failed to delete api keys:", err)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed to commit delete api keys tx:", err)
		return err
	}

	return nil
}
