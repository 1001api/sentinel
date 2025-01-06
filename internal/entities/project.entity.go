package entities

import "time"

type Project struct {
	ID          string
	Name        string
	Url         string
	Description string
	UserID      string
	CreatedAt   time.Time
}

type ProjectAggr struct {
	TotalEvents          int32
	TotalEventTypes      int32
	TotalUniqueUsers     int32
	TotalLocations       int32
	TotalUniquePageUrls  int32
	MostVisitedUrls      interface{}
	MostVisitedCountries interface{}
	MostVisitedCities    interface{}
	MostVisitedRegions   interface{}
	MostFiringElements   interface{}
	LastVisitedUsers     interface{}
	MostUsedBrowsers     interface{}
	MostFiredEventTypes  interface{}
	MostFiredEventLabels interface{}
	AggregatedAtStr      string
}
