package services

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
)

type DownloadService interface {
	DownloadEventData(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]string, [][]string, error)
}

type DownloadServiceImpl struct {
	Repo *gen.Queries
}

func (s *DownloadServiceImpl) DownloadEventData(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]string, [][]string, error) {
	header, err := s.Repo.GetEventTableHeaders(ctx)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	body, err := s.Repo.DownloadLastMonthData(ctx, gen.DownloadLastMonthDataParams{
		UserID:    userID,
		ProjectID: projectID,
	})
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	var result [][]string
	for _, row := range body {
		var ip string

		if row.IpAddr != nil {
			ip = row.IpAddr.String()
		} else {
			ip = ""
		}

		item := []string{
			row.ID.String(),
			row.EventType,
			row.EventLabel.String,
			row.PageUrl.String,
			row.ElementPath.String,
			row.ElementType.String,
			ip,
			row.UserAgent.String,
			row.BrowserName.String,
			row.Country.String,
			row.Region.String,
			row.City.String,
			row.SessionID.String,
			row.DeviceType.String,
			fmt.Sprintf("%d", row.TimeOnPage.Int32),
			row.ScreenResolution.String,
			row.FiredAt.String(),
			row.ReceivedAt.String(),
			row.UserID.String(),
			row.ProjectID.String(),
		}

		result = append(result, item)
	}

	return header, result, nil
}
