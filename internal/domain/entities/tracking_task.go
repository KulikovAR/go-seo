package entities

import (
	"time"
)

type TrackingTaskStatus string

const (
	TaskStatusPending   TrackingTaskStatus = "pending"
	TaskStatusRunning   TrackingTaskStatus = "running"
	TaskStatusCompleted TrackingTaskStatus = "completed"
	TaskStatusFailed    TrackingTaskStatus = "failed"
	TaskStatusCancelled TrackingTaskStatus = "cancelled"
)

type TrackingJob struct {
	ID             string             `json:"id"`
	SiteID         int                `json:"site_id"`
	Source         string             `json:"source"`
	Status         TrackingTaskStatus `json:"status"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
	CompletedAt    *time.Time         `json:"completed_at,omitempty"`
	TotalTasks     int                `json:"total_tasks"`
	CompletedTasks int                `json:"completed_tasks"`
	FailedTasks    int                `json:"failed_tasks"`
	Error          string             `json:"error,omitempty"`
}

type TrackingTask struct {
	ID          string             `json:"id"`
	JobID       string             `json:"job_id"`
	KeywordID   int                `json:"keyword_id"`
	SiteID      int                `json:"site_id"`
	Source      string             `json:"source"`
	Status      TrackingTaskStatus `json:"status"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
	CompletedAt *time.Time         `json:"completed_at,omitempty"`
	RetryCount  int                `json:"retry_count"`
	MaxRetries  int                `json:"max_retries"`
	Error       string             `json:"error,omitempty"`
	// Task-specific parameters
	Device     string `json:"device,omitempty"`
	OS         string `json:"os,omitempty"`
	Ads        bool   `json:"ads"`
	Country    string `json:"country,omitempty"`
	Lang       string `json:"lang,omitempty"`
	Pages      int    `json:"pages"`
	Subdomains bool   `json:"subdomains"`
	// External service parameters
	XMLUserID  string `json:"xml_user_id,omitempty"`
	XMLAPIKey  string `json:"xml_api_key,omitempty"`
	XMLBaseURL string `json:"xml_base_url,omitempty"`
	// Source-specific parameters
	TBS               string `json:"tbs,omitempty"`
	Filter            int    `json:"filter"`
	Highlights        int    `json:"highlights"`
	NFPR              int    `json:"nfpr"`
	Loc               int    `json:"loc"`
	AI                int    `json:"ai"`
	Raw               string `json:"raw,omitempty"`
	GroupBy           int    `json:"groupby"`
	Within            int    `json:"within"`
	LR                int    `json:"lr"`
	Domain            int    `json:"domain"`
	InIndex           int    `json:"inindex"`
	Strict            int    `json:"strict"`
	Organic           bool   `json:"organic"`
	Regions           *int   `json:"regions,omitempty"`
	FilterGroupID     *int   `json:"filter_group_id,omitempty"`
	WordstatQueryType string `json:"wordstat_query_type,omitempty"`
}

type TrackingResult struct {
	TaskID    string    `json:"task_id"`
	JobID     string    `json:"job_id"`
	KeywordID int       `json:"keyword_id"`
	SiteID    int       `json:"site_id"`
	Source    string    `json:"source"`
	Rank      int       `json:"rank"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Device    string    `json:"device"`
	OS        string    `json:"os"`
	Ads       bool      `json:"ads"`
	Country   string    `json:"country"`
	Lang      string    `json:"lang"`
	Pages     int       `json:"pages"`
	Date      time.Time `json:"date"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
}
