package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
)

type KeyRepo interface {
	CreateAPIKey(ctx context.Context, input *gen.CreateAPIKeyParams) (gen.CreateAPIKeyRow, error)
	GetAllPrivateKeys(ctx context.Context, userID uuid.UUID) ([]gen.FindAllAPIKeysRow, error)
	DeletePrivateKey(ctx context.Context, input *gen.DeleteAPIKeyParams) error
}

type KeyRepoImpl struct {
	Repo *gen.Queries
}

func InitKeyRepo(repo *gen.Queries) KeyRepoImpl {
	return KeyRepoImpl{
		Repo: repo,
	}
}

func (r *KeyRepoImpl) CreateAPIKey(ctx context.Context, input *gen.CreateAPIKeyParams) (gen.CreateAPIKeyRow, error) {
	return r.Repo.CreateAPIKey(ctx, *input)
}

func (r *KeyRepoImpl) GetAllPrivateKeys(ctx context.Context, userID uuid.UUID) ([]gen.FindAllAPIKeysRow, error) {
	return r.Repo.FindAllAPIKeys(ctx, userID)
}

func (r *KeyRepoImpl) DeletePrivateKey(ctx context.Context, input *gen.DeleteAPIKeyParams) error {
	return r.Repo.DeleteAPIKey(ctx, *input)
}
