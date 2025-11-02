package models

import (
	"time"
)

type TrackingJob struct {
	ID             string    `gorm:"primaryKey;type:varchar(50)"`
	SiteID         int       `gorm:"not null;index"`
	Source         string    `gorm:"not null;type:varchar(20)"`
	Status         string    `gorm:"not null;type:varchar(20);index"`
	CreatedAt      time.Time `gorm:"not null"`
	UpdatedAt      time.Time `gorm:"not null"`
	CompletedAt    *time.Time
	TotalTasks     int    `gorm:"not null;default:0"`
	CompletedTasks int    `gorm:"not null;default:0"`
	FailedTasks    int    `gorm:"not null;default:0"`
	Error          string `gorm:"type:text"`
}

func (TrackingJob) TableName() string {
	return "tracking_jobs"
}

type TrackingTask struct {
	ID          string    `gorm:"primaryKey;type:varchar(50)"`
	JobID       string    `gorm:"not null;type:varchar(50);index"`
	KeywordID   int       `gorm:"not null;index"`
	SiteID      int       `gorm:"not null;index"`
	Source      string    `gorm:"not null;type:varchar(20)"`
	Status      string    `gorm:"not null;type:varchar(20);index"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	CompletedAt *time.Time
	RetryCount  int    `gorm:"not null;default:0"`
	MaxRetries  int    `gorm:"not null;default:5"`
	Error       string `gorm:"type:text"`
	Device      string `gorm:"type:varchar(20)"`
	OS          string `gorm:"type:varchar(20)"`
	Ads         bool   `gorm:"default:false"`
	Country     string `gorm:"type:varchar(10)"`
	Lang        string `gorm:"type:varchar(10)"`
	Pages       int    `gorm:"default:0"`
	Subdomains  bool   `gorm:"default:false"`
	XMLUserID   string `gorm:"type:varchar(100)"`
	XMLAPIKey   string `gorm:"type:varchar(100)"`
	XMLBaseURL  string `gorm:"type:varchar(200)"`
	TBS         string `gorm:"type:varchar(50)"`
	Filter      int    `gorm:"default:0"`
	Highlights  int    `gorm:"default:0"`
	NFPR        int    `gorm:"default:0"`
	Loc         int    `gorm:"default:0"`
	AI          int    `gorm:"default:0"`
	Raw         string `gorm:"type:varchar(50)"`
	GroupBy     int    `gorm:"default:0"`
	Within      int    `gorm:"default:0"`
	LR          int    `gorm:"default:0"`
	Domain      string `gorm:"type:varchar(50)"`
	InIndex     int    `gorm:"default:0"`
	Strict      int    `gorm:"default:0"`
	Regions     *int   `gorm:"type:integer"`
}

func (TrackingTask) TableName() string {
	return "tracking_tasks"
}

type TrackingResult struct {
	ID        uint      `gorm:"primaryKey"`
	TaskID    string    `gorm:"not null;type:varchar(50);index"`
	JobID     string    `gorm:"not null;type:varchar(50);index"`
	KeywordID int       `gorm:"not null;index"`
	SiteID    int       `gorm:"not null;index"`
	Source    string    `gorm:"not null;type:varchar(20)"`
	Rank      int       `gorm:"not null"`
	URL       string    `gorm:"type:text"`
	Title     string    `gorm:"type:text"`
	Device    string    `gorm:"type:varchar(20)"`
	OS        string    `gorm:"type:varchar(20)"`
	Ads       bool      `gorm:"default:false"`
	Country   string    `gorm:"type:varchar(10)"`
	Lang      string    `gorm:"type:varchar(10)"`
	Pages     int       `gorm:"default:0"`
	Date      time.Time `gorm:"not null"`
	Success   bool      `gorm:"not null"`
	Error     string    `gorm:"type:text"`
}

func (TrackingResult) TableName() string {
	return "tracking_results"
}
