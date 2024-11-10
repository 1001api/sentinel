package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
)

type APIKeyRepo interface {
	CreateAPIKey(ctx context.Context, input *gen.CreateAPIKeyParams) (gen.CreateAPIKeyRow, error)
	GetAllPrivateKeys(ctx context.Context, userID uuid.UUID) ([]gen.FindAllAPIKeysRow, error)
	DeletePrivateKey(ctx context.Context, input *gen.DeleteAPIKeyParams) error
}

type APIKeyRepoImpl struct {
	Repo *gen.Queries
}

func InitAPIKeyRepo(repo *gen.Queries) APIKeyRepoImpl {
	return APIKeyRepoImpl{
		Repo: repo,
	}
}

func (r *APIKeyRepoImpl) CreateAPIKey(ctx context.Context, input *gen.CreateAPIKeyParams) (gen.CreateAPIKeyRow, error) {
	return r.Repo.CreateAPIKey(ctx, *input)
}

func (r *APIKeyRepoImpl) GetAllPrivateKeys(ctx context.Context, userID uuid.UUID) ([]gen.FindAllAPIKeysRow, error) {
	return r.Repo.FindAllAPIKeys(ctx, userID)
}

func (r *APIKeyRepoImpl) DeletePrivateKey(ctx context.Context, input *gen.DeleteAPIKeyParams) error {
	return r.Repo.DeleteAPIKey(ctx, *input)
}
