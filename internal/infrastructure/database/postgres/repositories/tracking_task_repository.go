package repositories

import (
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type TrackingTaskRepository struct {
	db *gorm.DB
}

func NewTrackingTaskRepository(db *gorm.DB) repositories.TrackingTaskRepository {
	return &TrackingTaskRepository{db: db}
}

func (r *TrackingTaskRepository) Create(task *entities.TrackingTask) error {
	model := &models.TrackingTask{
		ID:          task.ID,
		JobID:       task.JobID,
		KeywordID:   task.KeywordID,
		SiteID:      task.SiteID,
		Source:      task.Source,
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		CompletedAt: task.CompletedAt,
		RetryCount:  task.RetryCount,
		MaxRetries:  task.MaxRetries,
		Error:       task.Error,
		Device:      task.Device,
		OS:          task.OS,
		Ads:         task.Ads,
		Country:     task.Country,
		Lang:        task.Lang,
		Pages:       task.Pages,
		Subdomains:  task.Subdomains,
		XMLUserID:   task.XMLUserID,
		XMLAPIKey:   task.XMLAPIKey,
		XMLBaseURL:  task.XMLBaseURL,
		TBS:         task.TBS,
		Filter:      task.Filter,
		Highlights:  task.Highlights,
		NFPR:        task.NFPR,
		Loc:         task.Loc,
		AI:          task.AI,
		Raw:         task.Raw,
		GroupBy:     task.GroupBy,
		Within:      task.Within,
		LR:          task.LR,
		InIndex:     task.InIndex,
		Strict:      task.Strict,
		Regions:     task.Regions,
	}

	return r.db.Create(model).Error
}

func (r *TrackingTaskRepository) GetByID(id string) (*entities.TrackingTask, error) {
	var model models.TrackingTask
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(&model), nil
}

func (r *TrackingTaskRepository) GetByJobID(jobID string) ([]*entities.TrackingTask, error) {
	var models []models.TrackingTask
	if err := r.db.Where("job_id = ?", jobID).Find(&models).Error; err != nil {
		return nil, err
	}

	var tasks []*entities.TrackingTask
	for _, model := range models {
		tasks = append(tasks, r.modelToEntity(&model))
	}

	return tasks, nil
}

func (r *TrackingTaskRepository) Update(task *entities.TrackingTask) error {
	model := &models.TrackingTask{
		ID:          task.ID,
		JobID:       task.JobID,
		KeywordID:   task.KeywordID,
		SiteID:      task.SiteID,
		Source:      task.Source,
		Status:      string(task.Status),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		CompletedAt: task.CompletedAt,
		RetryCount:  task.RetryCount,
		MaxRetries:  task.MaxRetries,
		Error:       task.Error,
		Device:      task.Device,
		OS:          task.OS,
		Ads:         task.Ads,
		Country:     task.Country,
		Lang:        task.Lang,
		Pages:       task.Pages,
		Subdomains:  task.Subdomains,
		XMLUserID:   task.XMLUserID,
		XMLAPIKey:   task.XMLAPIKey,
		XMLBaseURL:  task.XMLBaseURL,
		TBS:         task.TBS,
		Filter:      task.Filter,
		Highlights:  task.Highlights,
		NFPR:        task.NFPR,
		Loc:         task.Loc,
		AI:          task.AI,
		Raw:         task.Raw,
		GroupBy:     task.GroupBy,
		Within:      task.Within,
		LR:          task.LR,
		InIndex:     task.InIndex,
		Strict:      task.Strict,
		Regions:     task.Regions,
	}

	return r.db.Save(model).Error
}

func (r *TrackingTaskRepository) UpdateStatus(id string, status entities.TrackingTaskStatus) error {
	return r.db.Model(&models.TrackingTask{}).
		Where("id = ?", id).
		Update("status", string(status)).Error
}

func (r *TrackingTaskRepository) UpdateRetryCount(id string, retryCount int) error {
	return r.db.Model(&models.TrackingTask{}).
		Where("id = ?", id).
		Update("retry_count", retryCount).Error
}

func (r *TrackingTaskRepository) GetPendingTasks(limit int) ([]*entities.TrackingTask, error) {
	var models []models.TrackingTask
	if err := r.db.Where("status = ?", string(entities.TaskStatusPending)).
		Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	var tasks []*entities.TrackingTask
	for _, model := range models {
		tasks = append(tasks, r.modelToEntity(&model))
	}

	return tasks, nil
}

func (r *TrackingTaskRepository) GetFailedTasks(limit int) ([]*entities.TrackingTask, error) {
	var models []models.TrackingTask
	if err := r.db.Where("status = ? AND retry_count < max_retries", string(entities.TaskStatusFailed)).
		Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}

	var tasks []*entities.TrackingTask
	for _, model := range models {
		tasks = append(tasks, r.modelToEntity(&model))
	}

	return tasks, nil
}

func (r *TrackingTaskRepository) Delete(id string) error {
	return r.db.Delete(&models.TrackingTask{}, "id = ?", id).Error
}

func (r *TrackingTaskRepository) DeleteByJobID(jobID string) error {
	return r.db.Delete(&models.TrackingTask{}, "job_id = ?", jobID).Error
}

func (r *TrackingTaskRepository) modelToEntity(model *models.TrackingTask) *entities.TrackingTask {
	return &entities.TrackingTask{
		ID:          model.ID,
		JobID:       model.JobID,
		KeywordID:   model.KeywordID,
		SiteID:      model.SiteID,
		Source:      model.Source,
		Status:      entities.TrackingTaskStatus(model.Status),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		CompletedAt: model.CompletedAt,
		RetryCount:  model.RetryCount,
		MaxRetries:  model.MaxRetries,
		Error:       model.Error,
		Device:      model.Device,
		OS:          model.OS,
		Ads:         model.Ads,
		Country:     model.Country,
		Lang:        model.Lang,
		Pages:       model.Pages,
		Subdomains:  model.Subdomains,
		XMLUserID:   model.XMLUserID,
		XMLAPIKey:   model.XMLAPIKey,
		XMLBaseURL:  model.XMLBaseURL,
		TBS:         model.TBS,
		Filter:      model.Filter,
		Highlights:  model.Highlights,
		NFPR:        model.NFPR,
		Loc:         model.Loc,
		AI:          model.AI,
		Raw:         model.Raw,
		GroupBy:     model.GroupBy,
		Within:      model.Within,
		LR:          model.LR,
		InIndex:     model.InIndex,
		Strict:      model.Strict,
		Regions:     model.Regions,
	}
}
