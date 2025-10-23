package repositories

import (
	"time"

	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type TrackingJobRepository struct {
	db *gorm.DB
}

func NewTrackingJobRepository(db *gorm.DB) repositories.TrackingJobRepository {
	return &TrackingJobRepository{db: db}
}

func (r *TrackingJobRepository) Create(job *entities.TrackingJob) error {
	model := &models.TrackingJob{
		ID:             job.ID,
		SiteID:         job.SiteID,
		Source:         job.Source,
		Status:         string(job.Status),
		CreatedAt:      job.CreatedAt,
		UpdatedAt:      job.UpdatedAt,
		CompletedAt:    job.CompletedAt,
		TotalTasks:     job.TotalTasks,
		CompletedTasks: job.CompletedTasks,
		FailedTasks:    job.FailedTasks,
		Error:          job.Error,
	}

	return r.db.Create(model).Error
}

func (r *TrackingJobRepository) GetByID(id string) (*entities.TrackingJob, error) {
	var model models.TrackingJob
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	return &entities.TrackingJob{
		ID:             model.ID,
		SiteID:         model.SiteID,
		Source:         model.Source,
		Status:         entities.TrackingTaskStatus(model.Status),
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		CompletedAt:    model.CompletedAt,
		TotalTasks:     model.TotalTasks,
		CompletedTasks: model.CompletedTasks,
		FailedTasks:    model.FailedTasks,
		Error:          model.Error,
	}, nil
}

func (r *TrackingJobRepository) Update(job *entities.TrackingJob) error {
	model := &models.TrackingJob{
		ID:             job.ID,
		SiteID:         job.SiteID,
		Source:         job.Source,
		Status:         string(job.Status),
		CreatedAt:      job.CreatedAt,
		UpdatedAt:      job.UpdatedAt,
		CompletedAt:    job.CompletedAt,
		TotalTasks:     job.TotalTasks,
		CompletedTasks: job.CompletedTasks,
		FailedTasks:    job.FailedTasks,
		Error:          job.Error,
	}

	return r.db.Save(model).Error
}

func (r *TrackingJobRepository) UpdateStatus(id string, status entities.TrackingTaskStatus) error {
	return r.db.Model(&models.TrackingJob{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     string(status),
			"updated_at": time.Now(),
		}).Error
}

func (r *TrackingJobRepository) UpdateProgress(id string, completed, failed int) error {
	return r.db.Model(&models.TrackingJob{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"completed_tasks": completed,
			"failed_tasks":    failed,
			"updated_at":      time.Now(),
		}).Error
}

func (r *TrackingJobRepository) GetBySiteID(siteID int) ([]*entities.TrackingJob, error) {
	var models []models.TrackingJob
	if err := r.db.Where("site_id = ?", siteID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	var jobs []*entities.TrackingJob
	for _, model := range models {
		jobs = append(jobs, &entities.TrackingJob{
			ID:             model.ID,
			SiteID:         model.SiteID,
			Source:         model.Source,
			Status:         entities.TrackingTaskStatus(model.Status),
			CreatedAt:      model.CreatedAt,
			UpdatedAt:      model.UpdatedAt,
			CompletedAt:    model.CompletedAt,
			TotalTasks:     model.TotalTasks,
			CompletedTasks: model.CompletedTasks,
			FailedTasks:    model.FailedTasks,
			Error:          model.Error,
		})
	}

	return jobs, nil
}

func (r *TrackingJobRepository) GetByStatus(status entities.TrackingTaskStatus) ([]*entities.TrackingJob, error) {
	var models []models.TrackingJob
	if err := r.db.Where("status = ?", string(status)).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	var jobs []*entities.TrackingJob
	for _, model := range models {
		jobs = append(jobs, &entities.TrackingJob{
			ID:             model.ID,
			SiteID:         model.SiteID,
			Source:         model.Source,
			Status:         entities.TrackingTaskStatus(model.Status),
			CreatedAt:      model.CreatedAt,
			UpdatedAt:      model.UpdatedAt,
			CompletedAt:    model.CompletedAt,
			TotalTasks:     model.TotalTasks,
			CompletedTasks: model.CompletedTasks,
			FailedTasks:    model.FailedTasks,
			Error:          model.Error,
		})
	}

	return jobs, nil
}

func (r *TrackingJobRepository) Delete(id string) error {
	return r.db.Delete(&models.TrackingJob{}, "id = ?", id).Error
}
