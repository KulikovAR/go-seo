package usecases

import (
	"go-seo/internal/delivery/http/dto"
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
)

type TrackingJobUseCase struct {
	trackingJobRepo repositories.TrackingJobRepository
}

func NewTrackingJobUseCase(trackingJobRepo repositories.TrackingJobRepository) *TrackingJobUseCase {
	return &TrackingJobUseCase{
		trackingJobRepo: trackingJobRepo,
	}
}

func (uc *TrackingJobUseCase) GetJobsWithPagination(req *dto.TrackingJobsRequest) (*dto.TrackingJobsResponse, error) {
	// Устанавливаем значения по умолчанию
	page := req.Page
	if page <= 0 {
		page = 1
	}
	perPage := req.PerPage
	if perPage <= 0 {
		perPage = 20
	}
	if perPage > 100 {
		perPage = 100
	}

	// Конвертируем статус в entity
	var status *entities.TrackingTaskStatus
	if req.Status != nil {
		s := entities.TrackingTaskStatus(*req.Status)
		status = &s
	}

	// Получаем данные из репозитория
	jobs, total, err := uc.trackingJobRepo.GetJobsWithPagination(page, perPage, req.SiteID, status)
	if err != nil {
		return nil, err
	}

	// Конвертируем в DTO
	var jobItems []dto.TrackingJobItem
	for _, job := range jobs {
		progress := 0.0
		if job.TotalTasks > 0 {
			progress = float64(job.CompletedTasks) / float64(job.TotalTasks) * 100
		}

		jobItems = append(jobItems, dto.TrackingJobItem{
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
			Progress:       progress,
		})
	}

	// Вычисляем пагинацию
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	from := (page-1)*perPage + 1
	to := page * perPage
	if to > int(total) {
		to = int(total)
	}
	if from > int(total) {
		from = 0
	}

	return &dto.TrackingJobsResponse{
		Data: jobItems,
		Pagination: dto.PaginationInfo{
			CurrentPage: page,
			PerPage:     perPage,
			Total:       int(total),
			LastPage:    lastPage,
			From:        from,
			To:          to,
			HasMore:     page < lastPage,
		},
		Meta: dto.MetaInfo{
			QueryTimeMs: 0, // Можно добавить измерение времени запроса
			Cached:      false,
		},
	}, nil
}
