package entities

import (
	"database/sql"
	"net"
	"time"
)

type Event struct {
	ID               string         `db:"id"`
	EventType        string         `json:"event_type" db:"event_type"`
	EventLabel       sql.NullString `json:"event_label,omitempty" db:"event_label"`
	PageURL          sql.NullString `json:"page_url,omitempty" db:"page_url"`
	ElementPath      sql.NullString `json:"element_path,omitempty" db:"element_path"`
	ElementType      sql.NullString `json:"element_type,omitempty" db:"element_type"`
	IPAddr           net.IP         `json:"ip_addr,omitempty" db:"ip_addr"`
	UserAgent        sql.NullString `json:"user_agent,omitempty" db:"user_agent"`
	BrowserName      sql.NullString `json:"browser_name,omitempty" db:"browser_name"`
	Country          sql.NullString `json:"country,omitempty" db:"country"`
	Region           sql.NullString `json:"region,omitempty" db:"region"`
	City             sql.NullString `json:"city,omitempty" db:"city"`
	SessionID        sql.NullString `json:"session_id,omitempty" db:"session_id"`
	DeviceType       sql.NullString `json:"device_type,omitempty" db:"device_type"`
	TimeOnPage       sql.NullInt32  `json:"time_on_page,omitempty" db:"time_on_page"`
	ScreenResolution sql.NullString `json:"screen_resolution,omitempty" db:"screen_resolution"`
	FiredAt          time.Time      `json:"fired_at" db:"fired_at"`
	ReceivedAt       time.Time      `json:"received_at" db:"received_at"`
	UserID           string         `json:"user_id" db:"user_id"`
	ProjectID        int            `json:"project_id" db:"project_id"`

	ProjectName string
}

type EventSummary struct {
	TotalEvents         int            `db:"total_events"`
	TotalUniqueUsers    int            `db:"total_unique_users"`
	TotalEventType      int            `db:"total_event_type"`
	TotalCountryVisited int            `db:"total_country_visited"`
	MostVisitedURL      sql.NullString `db:"most_visited_url"`
	MostCountryVisited  sql.NullString `db:"most_country_visited"`
}

type EventTextTotal struct {
	Name  string `db:"name"`
	Total int    `db:"total"`
}

type EventLastUser struct {
	IP        net.IP    `db:"ip"`
	Timestamp time.Time `db:"timestamp"`
}

type EventDetail struct {
	TotalEvents         int `db:"total_events"`
	TotalEventType      int `db:"total_event_type"`
	TotalUniqueUsers    int `db:"total_unique_users"`
	TotalCountryVisited int `db:"total_country_visited"`
	TotalPageURL        int `db:"total_page_url"`
	MostVisitedURLs     []EventTextTotal
	MostCountryVisited  []EventTextTotal
	MostCitiesVisited   []EventTextTotal
	MostElementsFired   []EventTextTotal
	LastVisitedUsers    []EventLastUser
	MostUsedBrowsers    []EventTextTotal
	MostFiredEventType  []EventTextTotal
	MostFiredEventLabel []EventTextTotal
}
