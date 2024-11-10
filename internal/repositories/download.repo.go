package repositories

import (
	"context"

	"github.com/hubkudev/sentinel/gen"
)

type DownloadRepo interface {
	GetEventTableHeaders(ctx context.Context) ([]string, error)
	DownloadLastMonthData(ctx context.Context, input *gen.DownloadLastMonthDataParams) ([]gen.Event, error)
}

type DownloadRepoImpl struct {
	Repo *gen.Queries
}

func InitDownloadRepo(repo *gen.Queries) DownloadRepoImpl {
	return DownloadRepoImpl{
		Repo: repo,
	}
}

func (r *DownloadRepoImpl) GetEventTableHeaders(ctx context.Context) ([]string, error) {
	return r.Repo.GetEventTableHeaders(ctx)
}

func (r *DownloadRepoImpl) DownloadLastMonthData(ctx context.Context, input *gen.DownloadLastMonthDataParams) ([]gen.Event, error) {
	return r.Repo.DownloadLastMonthData(ctx, *input)
}
