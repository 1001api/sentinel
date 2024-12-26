package configs

import (
	"fmt"

	"github.com/google/uuid"
)

func CACHE_LIVE_EVENTS(id uuid.UUID) string {
	return fmt.Sprintf("cache:live-events/%s", id)
}

func CACHE_LIVE_EVENT(id uuid.UUID, projectID uuid.UUID, strategy string) string {
	return fmt.Sprintf("cache:live-event/%s/%s?strategy=%s", id, projectID, strategy)
}

func CACHE_LIVE_EVENT_SUMMARY(userID uuid.UUID, projectID uuid.UUID) string {
	return fmt.Sprintf("cache:live-event-summary/%s/%s", userID, projectID)
}

func CACHE_LIVE_EVENT_DETAIL(userID uuid.UUID, projectID uuid.UUID) string {
	return fmt.Sprintf("cache:live-event-detail/%s/%s", userID, projectID)
}

func CACHE_LIVE_EVENT_DETAIL_SUMMARY(userID uuid.UUID, projectID uuid.UUID) string {
	return fmt.Sprintf("cache:live-event-detail-summary/%s/%s", userID, projectID)
}

func CACHE_JSON_WEEKLY_EVENT_CHART(userID uuid.UUID, projectID uuid.UUID) string {
	return fmt.Sprintf("cache:json-weekly-event-chart/%s/%s", userID, projectID)
}

func CACHE_JSON_EVENT_TYPE_CHART(userID uuid.UUID, projectID uuid.UUID) string {
	return fmt.Sprintf("cache:json-event-type-chart/%s/%s", userID, projectID)
}

func CACHE_JSON_EVENT_LABEL_CHART(userID uuid.UUID, projectID uuid.UUID) string {
	return fmt.Sprintf("cache:json-event-label-chart/%s/%s", userID, projectID)
}
