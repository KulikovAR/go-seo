package repositories

import (
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type TrackingResultRepository struct {
	db *gorm.DB
}

func NewTrackingResultRepository(db *gorm.DB) repositories.TrackingResultRepository {
	return &TrackingResultRepository{db: db}
}

func (r *TrackingResultRepository) Create(result *entities.TrackingResult) error {
	model := &models.TrackingResult{
		TaskID:    result.TaskID,
		JobID:     result.JobID,
		KeywordID: result.KeywordID,
		SiteID:    result.SiteID,
		Source:    result.Source,
		Rank:      result.Rank,
		URL:       result.URL,
		Title:     result.Title,
		Device:    result.Device,
		OS:        result.OS,
		Ads:       result.Ads,
		Country:   result.Country,
		Lang:      result.Lang,
		Pages:     result.Pages,
		Date:      result.Date,
		Success:   result.Success,
		Error:     result.Error,
	}

	return r.db.Create(model).Error
}

func (r *TrackingResultRepository) GetByTaskID(taskID string) (*entities.TrackingResult, error) {
	var model models.TrackingResult
	if err := r.db.Where("task_id = ?", taskID).First(&model).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(&model), nil
}

func (r *TrackingResultRepository) GetByJobID(jobID string) ([]*entities.TrackingResult, error) {
	var models []models.TrackingResult
	if err := r.db.Where("job_id = ?", jobID).Find(&models).Error; err != nil {
		return nil, err
	}

	var results []*entities.TrackingResult
	for _, model := range models {
		results = append(results, r.modelToEntity(&model))
	}

	return results, nil
}

func (r *TrackingResultRepository) GetBySiteID(siteID int) ([]*entities.TrackingResult, error) {
	var models []models.TrackingResult
	if err := r.db.Where("site_id = ?", siteID).Find(&models).Error; err != nil {
		return nil, err
	}

	var results []*entities.TrackingResult
	for _, model := range models {
		results = append(results, r.modelToEntity(&model))
	}

	return results, nil
}

func (r *TrackingResultRepository) Delete(id string) error {
	return r.db.Delete(&models.TrackingResult{}, "id = ?", id).Error
}

func (r *TrackingResultRepository) DeleteByJobID(jobID string) error {
	return r.db.Delete(&models.TrackingResult{}, "job_id = ?", jobID).Error
}

func (r *TrackingResultRepository) modelToEntity(model *models.TrackingResult) *entities.TrackingResult {
	return &entities.TrackingResult{
		TaskID:    model.TaskID,
		JobID:     model.JobID,
		KeywordID: model.KeywordID,
		SiteID:    model.SiteID,
		Source:    model.Source,
		Rank:      model.Rank,
		URL:       model.URL,
		Title:     model.Title,
		Device:    model.Device,
		OS:        model.OS,
		Ads:       model.Ads,
		Country:   model.Country,
		Lang:      model.Lang,
		Pages:     model.Pages,
		Date:      model.Date,
		Success:   model.Success,
		Error:     model.Error,
	}
}
