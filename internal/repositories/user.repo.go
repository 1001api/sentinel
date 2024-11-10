package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
)

type UserRepo interface {
	FindUserByEmail(ctx context.Context, email string) (gen.FindUserByEmailRow, error)
	FindUserByID(ctx context.Context, id uuid.UUID) (gen.FindUserByIDRow, error)
	FindUserByEmailWithHash(ctx context.Context, email string) (gen.FindUserByEmailWithHashRow, error)
	FindUserByPublicKey(ctx context.Context, key string) (gen.FindUserByPublicKeyRow, error)
	FindUserPublicKey(ctx context.Context, id uuid.UUID) (string, error)
	CheckAdminExist(ctx context.Context) (bool, error)
	CreateUser(ctx context.Context, input *gen.CreateUserParams) (gen.CreateUserRow, error)
}

type UserRepoImpl struct {
	Repo *gen.Queries
}

func InitUserRepo(repo *gen.Queries) UserRepoImpl {
	return UserRepoImpl{
		Repo: repo,
	}
}

func (r *UserRepoImpl) FindUserByEmail(ctx context.Context, email string) (gen.FindUserByEmailRow, error) {
	return r.Repo.FindUserByEmail(ctx, email)
}

func (r *UserRepoImpl) FindUserByID(ctx context.Context, id uuid.UUID) (gen.FindUserByIDRow, error) {
	return r.Repo.FindUserByID(ctx, id)
}

func (r *UserRepoImpl) FindUserByEmailWithHash(ctx context.Context, email string) (gen.FindUserByEmailWithHashRow, error) {
	return r.Repo.FindUserByEmailWithHash(ctx, email)
}

func (r *UserRepoImpl) FindUserByPublicKey(ctx context.Context, key string) (gen.FindUserByPublicKeyRow, error) {
	return r.Repo.FindUserByPublicKey(ctx, key)
}

func (r *UserRepoImpl) FindUserPublicKey(ctx context.Context, id uuid.UUID) (string, error) {
	return r.Repo.FindUserPublicKey(ctx, id)
}

func (r *UserRepoImpl) CheckAdminExist(ctx context.Context) (bool, error) {
	return r.Repo.CheckAdminExist(ctx)
}

func (r *UserRepoImpl) CreateUser(ctx context.Context, input *gen.CreateUserParams) (gen.CreateUserRow, error) {
	return r.Repo.CreateUser(ctx, *input)
}
