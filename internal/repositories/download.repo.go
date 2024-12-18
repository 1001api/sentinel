package repositories

import (
	"context"

	"github.com/hubkudev/sentinel/gen"
)

type DownloadRepo interface {
	GetEventTableHeaders(ctx context.Context) ([]string, error)
	DownloadIntervalData(ctx context.Context, input *gen.DownloadIntervalEventDataParams) ([]gen.Event, error)
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

func (r *DownloadRepoImpl) DownloadIntervalData(ctx context.Context, input *gen.DownloadIntervalEventDataParams) ([]gen.Event, error) {
	return r.Repo.DownloadIntervalEventData(ctx, *input)
}
