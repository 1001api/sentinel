package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
)

type SubService interface {
	CheckUserHasActiveSub(ctx context.Context, userID string) (bool, error)
}

type SubServiceImpl struct {
	Repo *gen.Queries
}

func (s *SubServiceImpl) CheckUserHasActiveSub(ctx context.Context, userID string) (bool, error) {
	userUUID := uuid.MustParse(userID)
	return s.Repo.CheckUserHasActiveSub(ctx, userUUID)
}
