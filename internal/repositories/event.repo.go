package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
)

type EventRepo interface {
	CheckProjectWithinUserID(ctx context.Context, input *gen.CheckProjectWithinUserIDParams) (bool, error)
	CreateEvent(ctx context.Context, input *gen.CreateEventParams) error
	GetLiveEvents(ctx context.Context, userID uuid.UUID) ([]gen.GetLiveEventsRow, error)
	GetLiveEventDetail(ctx context.Context, input *gen.GetLiveEventsDetailParams) ([]gen.GetLiveEventsDetailRow, error)
	GetEventSummary(ctx context.Context, input *gen.GetEventSummaryParams) (gen.GetEventSummaryRow, error)
	GetTotalEventSummary(ctx context.Context, input *gen.GetTotalEventSummaryParams) (gen.GetTotalEventSummaryRow, error)
	GetEventDetailSummary(ctx context.Context, input *gen.GetEventDetailSummaryParams) ([]gen.GetEventDetailSummaryRow, error)
	GetWeeklyEvents(ctx context.Context, input *gen.GetWeeklyEventsParams) ([]gen.GetWeeklyEventsRow, error)
	GetWeeklyEventsTotal(ctx context.Context, input *gen.GetWeeklyEventsTotalParams) (int64, error)
	GetPercentageEventsType(ctx context.Context, input *gen.GetPercentageEventsTypeParams) ([]gen.GetPercentageEventsTypeRow, error)
	GetPercentageEventsLabel(ctx context.Context, input *gen.GetPercentageEventsLabelParams) ([]gen.GetPercentageEventsLabelRow, error)
	CountUserMonthlyEvents(ctx context.Context, userID uuid.UUID) (int64, error)
}

type EventRepoImpl struct {
	Repo *gen.Queries
}

func InitEventRepo(repo *gen.Queries) EventRepoImpl {
	return EventRepoImpl{
		Repo: repo,
	}
}

func (r *EventRepoImpl) CheckProjectWithinUserID(ctx context.Context, input *gen.CheckProjectWithinUserIDParams) (bool, error) {
	return r.Repo.CheckProjectWithinUserID(ctx, *input)
}

func (r *EventRepoImpl) CreateEvent(ctx context.Context, input *gen.CreateEventParams) error {
	return r.Repo.CreateEvent(ctx, *input)
}

func (r *EventRepoImpl) GetLiveEvents(ctx context.Context, userID uuid.UUID) ([]gen.GetLiveEventsRow, error) {
	return r.Repo.GetLiveEvents(ctx, userID)
}

func (r *EventRepoImpl) GetLiveEventDetail(ctx context.Context, input *gen.GetLiveEventsDetailParams) ([]gen.GetLiveEventsDetailRow, error) {
	return r.Repo.GetLiveEventsDetail(ctx, *input)
}

func (r *EventRepoImpl) GetEventSummary(ctx context.Context, input *gen.GetEventSummaryParams) (gen.GetEventSummaryRow, error) {
	return r.Repo.GetEventSummary(ctx, *input)
}

func (r *EventRepoImpl) GetTotalEventSummary(ctx context.Context, input *gen.GetTotalEventSummaryParams) (gen.GetTotalEventSummaryRow, error) {
	return r.Repo.GetTotalEventSummary(ctx, *input)
}

func (r *EventRepoImpl) GetEventDetailSummary(ctx context.Context, input *gen.GetEventDetailSummaryParams) ([]gen.GetEventDetailSummaryRow, error) {
	return r.Repo.GetEventDetailSummary(ctx, *input)
}

func (r *EventRepoImpl) GetWeeklyEvents(ctx context.Context, input *gen.GetWeeklyEventsParams) ([]gen.GetWeeklyEventsRow, error) {
	return r.Repo.GetWeeklyEvents(ctx, *input)
}

func (r *EventRepoImpl) GetWeeklyEventsTotal(ctx context.Context, input *gen.GetWeeklyEventsTotalParams) (int64, error) {
	return r.Repo.GetWeeklyEventsTotal(ctx, *input)
}

func (r *EventRepoImpl) GetPercentageEventsType(
	ctx context.Context,
	input *gen.GetPercentageEventsTypeParams,
) ([]gen.GetPercentageEventsTypeRow, error) {
	return r.Repo.GetPercentageEventsType(ctx, *input)
}

func (r *EventRepoImpl) GetPercentageEventsLabel(
	ctx context.Context,
	input *gen.GetPercentageEventsLabelParams,
) ([]gen.GetPercentageEventsLabelRow, error) {
	return r.Repo.GetPercentageEventsLabel(ctx, *input)
}

func (r *EventRepoImpl) CountUserMonthlyEvents(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.Repo.CountUserMonthlyEvents(ctx, userID)
}
