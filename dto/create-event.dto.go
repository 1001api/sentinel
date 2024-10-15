package dto

type CreateEventInput struct {
	ProjectID        string `json:"ProjectID" validate:"required,uuid"`
	EventType        string `json:"EventType" validate:"required,max=100"`
	EventLabel       string `json:"EventLabel,omitempty" validate:"omitempty,max=100"`
	PageURL          string `json:"PageURL,omitempty" validate:"omitempty,url"`
	ElementPath      string `json:"ElementPath,omitempty" validate:"omitempty,max=255"`
	ElementType      string `json:"ElementType,omitempty" validate:"omitempty,max=255"`
	IPAddr           string `json:"IPAddr,omitempty" validate:"omitempty,ip_addr"`
	UserAgent        string `json:"UserAgent,omitempty" validate:"omitempty,max=255"`
	BrowserName      string `json:"BrowserName,omitempty" validate:"omitempty,max=100"`
	Country          string `json:"Country,omitempty" validate:"omitempty,max=100"`
	Region           string `json:"Region,omitempty" validate:"omitempty,max=100"`
	City             string `json:"City,omitempty" validate:"omitempty,max=100"`
	SessionID        string `json:"SessionID,omitempty" validate:"omitempty,max=100"`
	DeviceType       string `json:"DeviceType,omitempty" validate:"omitempty,max=100"`
	TimeOnPage       int    `json:"TimeOnPage,omitempty"`
	ScreenResolution string `json:"ScreenResolution,omitempty" validate:"omitempty,max=100"`
	FiredAt          string `json:"FiredAt" validate:"required,timestamp"`
}
