package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// JellyfinConfig holds the configuration for connecting to a Jellyfin server
type JellyfinConfig struct {
	BaseURL       string `json:"base_url" yaml:"base_url"`
	Token         string `json:"token" yaml:"token"`
	SkipSSLVerify bool   `json:"insecure" yaml:"insecure"`
}

// LogConfig holds the logging configuration
type LogConfig struct {
	Level string `json:"level" yaml:"level"`
}

// Config represents the application configuration
type Config struct {
	Jellyfin JellyfinConfig `json:"api" yaml:"api"`
	Logging  LogConfig      `json:"logging" yaml:"logging"`
}

// Session represents a Jellyfin user session
type Session struct {
	DeviceName      string    `json:"DeviceName"`
	UserName        string    `json:"UserName"`
	LastActivityUTC time.Time `json:"LastActivityDate"`
	ClientName      string    `json:"Client"`
	ID              string    `json:"Id"`
}

// LibraryFolder represents a Jellyfin library folder
type LibraryFolder struct {
	Name               string                 `json:"Name"`
	CollectionType     string                 `json:"CollectionType"`
	RefreshStatus      string                 `json:"RefreshStatus"`
	ItemID             string                 `json:"ItemId"`
	PrimaryImageTag    string                 `json:"PrimaryImageTag"`
	AdditionalMetadata map[string]interface{} `json:"-"`
}

// ActivityLogItem represents a single activity log entry
type ActivityLogItem struct {
	ID               int64     `json:"Id"`
	Name             string    `json:"Name"`
	Overview         string    `json:"Overview"`
	ShortOverview    string    `json:"ShortOverview"`
	Type             string    `json:"Type"`
	ItemID           string    `json:"ItemId"`
	DateCreatedUTC   time.Time `json:"Date"`
	UserID           string    `json:"UserId"`
	UserPrimaryImage string    `json:"UserPrimaryImageTag"`
	Severity         string    `json:"Severity"`
}

// ActivityLog represents a collection of activity log entries
type ActivityLog struct {
	Items      []ActivityLogItem `json:"Items"`
	TotalCount int               `json:"TotalRecordCount"`
	StartIndex int               `json:"StartIndex"`
}

// SearchHint represents a search hint returned from Jellyfin
type SearchHint struct {
	Name        string `json:"Name"`
	ID          string `json:"ItemId"`
	Type        string `json:"Type"`
	MediaType   string `json:"MediaType"`
	PrimaryTag  string `json:"PrimaryImageTag"`
	SeriesName  string `json:"SeriesName"`
	EpisodeNum  int    `json:"IndexNumber"`
	SeasonNum   int    `json:"ParentIndexNumber"`
	ProductYear int    `json:"ProductionYear"`
}

// SearchResponse represents a search response from Jellyfin
type SearchResponse struct {
	SearchHints []SearchHint `json:"SearchHints"`
	TotalHints  int          `json:"TotalRecordCount"`
}

// UnmarshalJSON is a custom unmarshaler for time.Time fields in Jellyfin API responses
func ParseJellyfinTime(data []byte) (time.Time, error) {
	var timeStr string
	if err := json.Unmarshal(data, &timeStr); err != nil {
		return time.Time{}, fmt.Errorf("failed to unmarshal time string: %w", err)
	}

	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	return t, nil
}
