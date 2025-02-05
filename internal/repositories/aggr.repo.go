package repositories

import (
	"context"

	"github.com/hubkudev/sentinel/gen"
)

type AggrRepo interface {
	GetTotalAggr(ctx context.Context, input *gen.GetTotalAggrParams) (gen.GetTotalAggrRow, error)
	GetBriefAggr(ctx context.Context, input *gen.GetBriefAggrParams) (gen.GetBriefAggrRow, error)
	GetDetailAggr(ctx context.Context, input *gen.GetDetailAggrParams) ([]gen.GetDetailAggrRow, error)
}

type AggrRepoImpl struct {
	Repo *gen.Queries
}

func InitAggrRepo(repo *gen.Queries) AggrRepoImpl {
	return AggrRepoImpl{
		Repo: repo,
	}
}

func (r *AggrRepoImpl) GetTotalAggr(ctx context.Context, input *gen.GetTotalAggrParams) (gen.GetTotalAggrRow, error) {
	return r.Repo.GetTotalAggr(ctx, *input)
}

func (r *AggrRepoImpl) GetBriefAggr(ctx context.Context, input *gen.GetBriefAggrParams) (gen.GetBriefAggrRow, error) {
	return r.Repo.GetBriefAggr(ctx, *input)
}

func (r *AggrRepoImpl) GetDetailAggr(ctx context.Context, input *gen.GetDetailAggrParams) ([]gen.GetDetailAggrRow, error) {
	return r.Repo.GetDetailAggr(ctx, *input)
}
