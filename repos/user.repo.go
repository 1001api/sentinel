package repositories

import (
	"context"
	"log"

	"github.com/hubkudev/sentinel/dto"
	"github.com/hubkudev/sentinel/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByID(ctx context.Context, userID string) (*entities.User, error)
	GetPublicKey(ctx context.Context, userID string) (string, error)
	FindByPublicKey(ctx context.Context, publicKey string) (*entities.User, error)
	CheckIDExist(ctx context.Context, userID string) (bool, error)
	CreateUser(ctx context.Context, input *dto.CreateUserInput) (*entities.User, error)
}

type UserRepoImpl struct {
	DB *pgxpool.Pool
}

func (r *UserRepoImpl) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User

	SQL := "SELECT id, fullname, email, profile_url FROM users WHERE email = $1"
	row := r.DB.QueryRow(ctx, SQL, email)
	if err := row.Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.ProfileURL,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoImpl) FindByID(ctx context.Context, userID string) (*entities.User, error) {
	var user entities.User

	SQL := "SELECT id, fullname, email, profile_url FROM users WHERE id = $1"
	row := r.DB.QueryRow(ctx, SQL, userID)
	if err := row.Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.ProfileURL,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoImpl) GetPublicKey(ctx context.Context, userID string) (string, error) {
	var key string

	SQL := "SELECT public_key FROM users WHERE id = $1"
	row := r.DB.QueryRow(ctx, SQL, userID)
	if err := row.Scan(
		&key,
	); err != nil {
		return "", err
	}

	return key, nil
}

func (r *UserRepoImpl) FindByPublicKey(ctx context.Context, publicKey string) (*entities.User, error) {
	var user entities.User

	SQL := "SELECT id, fullname, email, profile_url FROM users WHERE id = $1"
	row := r.DB.QueryRow(ctx, SQL, publicKey)
	if err := row.Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.ProfileURL,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoImpl) CheckIDExist(ctx context.Context, userID string) (bool, error) {
	var exist bool

	SQL := "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)"
	row := r.DB.QueryRow(ctx, SQL, userID)
	if err := row.Scan(&exist); err != nil {
		return false, err
	}

	return exist, nil
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, input *dto.CreateUserInput) (*entities.User, error) {
	var user entities.User

	SQL := "INSERT INTO users(fullname, email, oauth_provider, oauth_id, profile_url, public_key) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, fullname, email, profile_url"

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

	row := tx.QueryRow(ctx, SQL, input.Fullname, input.Email, input.OAuthProvider, input.OAuthID, input.ProfileURL, input.PublicKey)
	if err := row.Scan(
		&user.ID,
		&user.Fullname,
		&user.Email,
		&user.ProfileURL,
	); err != nil {
		log.Println("Failed creating user:", err)
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Failed committing tx:", err)
		return nil, err
	}

	return &user, nil
}
