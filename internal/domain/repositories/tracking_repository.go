package repositories

import (
	"go-seo/internal/domain/entities"
)

// TrackingJobRepository defines the interface for tracking job operations
type TrackingJobRepository interface {
	Create(job *entities.TrackingJob) error
	GetByID(id string) (*entities.TrackingJob, error)
	Update(job *entities.TrackingJob) error
	UpdateStatus(id string, status entities.TrackingTaskStatus) error
	UpdateProgress(id string, completed, failed int) error
	GetBySiteID(siteID int) ([]*entities.TrackingJob, error)
	GetByStatus(status entities.TrackingTaskStatus) ([]*entities.TrackingJob, error)
	GetJobsWithPagination(page, perPage int, siteID *int, status *entities.TrackingTaskStatus) ([]*entities.TrackingJob, int64, error)
	Delete(id string) error
}

// TrackingTaskRepository defines the interface for tracking task operations
type TrackingTaskRepository interface {
	Create(task *entities.TrackingTask) error
	GetByID(id string) (*entities.TrackingTask, error)
	GetByJobID(jobID string) ([]*entities.TrackingTask, error)
	Update(task *entities.TrackingTask) error
	UpdateStatus(id string, status entities.TrackingTaskStatus) error
	UpdateRetryCount(id string, retryCount int) error
	GetPendingTasks(limit int) ([]*entities.TrackingTask, error)
	GetFailedTasks(limit int) ([]*entities.TrackingTask, error)
	Delete(id string) error
	DeleteByJobID(jobID string) error
}

// TrackingResultRepository defines the interface for tracking result operations
type TrackingResultRepository interface {
	Create(result *entities.TrackingResult) error
	GetByTaskID(taskID string) (*entities.TrackingResult, error)
	GetByJobID(jobID string) ([]*entities.TrackingResult, error)
	GetBySiteID(siteID int) ([]*entities.TrackingResult, error)
	Delete(id string) error
	DeleteByJobID(jobID string) error
}
