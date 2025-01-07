package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hubkudev/sentinel/gen"
	"github.com/hubkudev/sentinel/internal/entities"
	"github.com/hubkudev/sentinel/internal/repositories"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectService interface {
	CreateProject(ctx context.Context, name string, desc string, url string, userID uuid.UUID) (*gen.CreateProjectRow, error)
	UpdateProject(ctx context.Context, name string, desc string, url string, projectID uuid.UUID, userID uuid.UUID) error
	GetProjectByID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*gen.FindProjectByIDRow, error)
	GetAllProjects(ctx context.Context, userID uuid.UUID) ([]gen.FindAllProjectsRow, error)
	GetProjectCount(ctx context.Context, userID uuid.UUID) (int64, error)
	DeleteProject(ctx context.Context, userID uuid.UUID, projectID uuid.UUID) error
	CountProjectSize(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (int64, error)
	LastProjectDataReceived(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*time.Time, error)
	SaveProjectSummary(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error
	FindProjectSummary(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, limit int32) ([]entities.ProjectAggr, error)
}

type ProjectServiceImpl struct {
	EventService EventService
	UtilService  UtilService
	Repo         repositories.ProjectRepo
}

func InitProjectService(repo repositories.ProjectRepo, eventService EventService, utilService UtilService) ProjectServiceImpl {
	return ProjectServiceImpl{
		Repo:         repo,
		EventService: eventService,
		UtilService:  utilService,
	}
}

func (s *ProjectServiceImpl) CreateProject(ctx context.Context, name string, desc string, url string, userID uuid.UUID) (*gen.CreateProjectRow, error) {
	input := gen.CreateProjectParams{
		Name: name,
		Description: pgtype.Text{
			String: desc,
			Valid:  desc != "",
		},
		Url: pgtype.Text{
			String: url,
			Valid:  url != "",
		},
		UserID: userID,
		CreatedAt: pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		},
	}

	key, err := s.Repo.Create(ctx, &input)
	if err != nil {
		return nil, err
	}

	return &key, nil
}

func (s *ProjectServiceImpl) UpdateProject(ctx context.Context, name string, desc string, url string, projectID uuid.UUID, userID uuid.UUID) error {
	input := gen.UpdateProjectParams{
		Name: name,
		Description: pgtype.Text{
			String: desc,
			Valid:  desc != "",
		},
		Url: pgtype.Text{
			String: url,
			Valid:  url != "",
		},
		ID:     projectID,
		UserID: userID,
	}
	return s.Repo.Update(ctx, &input)
}

func (s *ProjectServiceImpl) GetProjectByID(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*gen.FindProjectByIDRow, error) {
	input := gen.FindProjectByIDParams{
		ID:     projectID,
		UserID: userID,
	}

	row, err := s.Repo.FindByID(ctx, &input)
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (s *ProjectServiceImpl) GetAllProjects(ctx context.Context, userID uuid.UUID) ([]gen.FindAllProjectsRow, error) {
	return s.Repo.FindAll(ctx, userID)
}

func (s *ProjectServiceImpl) GetProjectCount(ctx context.Context, userID uuid.UUID) (int64, error) {
	return s.Repo.Count(ctx, userID)
}

func (s *ProjectServiceImpl) DeleteProject(ctx context.Context, userID uuid.UUID, projectID uuid.UUID) error {
	input := gen.DeleteProjectParams{
		UserID: userID,
		ID:     projectID,
	}
	return s.Repo.Delete(ctx, &input)
}

func (s *ProjectServiceImpl) CountProjectSize(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (int64, error) {
	input := gen.CountProjectSizeParams{
		UserID:    userID,
		ProjectID: projectID,
	}

	size, err := s.Repo.CountSize(ctx, &input)
	if err != nil {
		return -1, err
	}
	return size, err
}

func (s *ProjectServiceImpl) LastProjectDataReceived(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*time.Time, error) {
	input := gen.LastProjectDataReceivedParams{
		UserID:    userID,
		ProjectID: projectID,
	}

	lastTime, err := s.Repo.LastDataReceived(ctx, &input)
	if err != nil {
		return nil, err
	}
	return &lastTime, err
}

func (s *ProjectServiceImpl) SaveProjectSummary(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) error {
	// check project aggregation eligibility.
	eligibleStatus, err := s.Repo.CheckProjectAggrEligibility(ctx, projectID)
	if err != nil {
		return err
	}

	// if eligibleStatus is less than 1, meaning there is no new event after last
	// project aggregation, just return and bypass this saving, dont create a new one.
	if eligibleStatus < 1 {
		return nil
	}

	summary, err := s.EventService.GetEventDetailSummary(ctx, projectID, userID)
	if err != nil {
		return err
	}

	urlJson, err := json.Marshal(summary.MostVisitedURLs)
	if err != nil {
		return err
	}

	citiesJson, err := json.Marshal(summary.MostCitiesVisited)
	if err != nil {
		return err
	}

	countryJson, err := json.Marshal(summary.MostCountryVisited)
	if err != nil {
		return err
	}

	elementsJson, err := json.Marshal(summary.MostElementsFired)
	if err != nil {
		return err
	}

	lastUsersJson, err := json.Marshal(summary.LastVisitedUsers)
	if err != nil {
		return err
	}

	browsersJson, err := json.Marshal(summary.MostUsedBrowsers)
	if err != nil {
		return err
	}

	eventTypesJson, err := json.Marshal(summary.MostFiredEventType)
	if err != nil {
		return err
	}

	eventLabelsJson, err := json.Marshal(summary.MostFiredEventLabel)
	if err != nil {
		return err
	}

	now := time.Now()

	input := gen.CreateProjectAggrParams{
		ProjectID:            projectID,
		UserID:               userID,
		TotalEvents:          int32(summary.TotalEvents),
		TotalEventTypes:      int32(summary.TotalEventType),
		TotalUniqueUsers:     int32(summary.TotalUniqueUsers),
		TotalLocations:       int32(summary.TotalCountryVisited),
		TotalUniquePageUrls:  int32(summary.TotalPageURL),
		MostVisitedUrls:      urlJson,
		MostVisitedCountries: citiesJson,
		MostVisitedCities:    citiesJson,
		MostVisitedRegions:   countryJson,
		MostFiringElements:   elementsJson,
		LastVisitedUsers:     lastUsersJson,
		MostUsedBrowsers:     browsersJson,
		MostFiredEventTypes:  eventTypesJson,
		MostFiredEventLabels: eventLabelsJson,
		AggregatedAt:         now,
		AggregatedAtStr:      now.Format("02/01/2006"),
	}

	if err := s.Repo.CreateProjectAggr(ctx, &input); err != nil {
		return err
	}

	return nil
}

func (s *ProjectServiceImpl) FindProjectSummary(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, limit int32) ([]entities.ProjectAggr, error) {
	// sync project summary first.
	// sync-ing project summary will only happen if the event is present
	// after the last project summarization. That means, if the event has not changed from the last
	// summarization, this function will insert nothing to db.
	if err := s.SaveProjectSummary(ctx, projectID, userID); err != nil {
		return nil, err
	}

	summary, err := s.Repo.FindProjectAggr(ctx, &gen.FindProjectAggrParams{
		ProjectID: projectID,
		UserID:    userID,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	var results []entities.ProjectAggr

	for _, v := range summary {
		var result entities.ProjectAggr
		result.TotalEvents = v.TotalEvents
		result.TotalEventTypes = v.TotalEventTypes
		result.TotalUniqueUsers = v.TotalUniqueUsers
		result.TotalLocations = v.TotalLocations
		result.TotalUniquePageUrls = v.TotalUniquePageUrls
		result.MostVisitedUrls = s.UtilService.ByteToJSON(v.MostVisitedUrls)
		result.MostVisitedCountries = s.UtilService.ByteToJSON(v.MostVisitedCountries)
		result.MostVisitedCities = s.UtilService.ByteToJSON(v.MostVisitedCities)
		result.MostVisitedRegions = s.UtilService.ByteToJSON(v.MostVisitedRegions)
		result.MostFiringElements = s.UtilService.ByteToJSON(v.MostFiringElements)
		result.LastVisitedUsers = s.UtilService.ByteToJSON(v.LastVisitedUsers)
		result.MostUsedBrowsers = s.UtilService.ByteToJSON(v.MostUsedBrowsers)
		result.MostFiredEventTypes = s.UtilService.ByteToJSON(v.MostFiredEventTypes)
		result.MostFiredEventLabels = s.UtilService.ByteToJSON(v.MostFiredEventLabels)
		result.AggregatedAtStr = v.AggregatedAtStr
		results = append(results, result)
	}

	return results, err
}
