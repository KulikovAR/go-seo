package handlers

import (
	"net/http"
	"time"

	"go-seo/internal/delivery/http/dto"
	"go-seo/internal/usecases"
	"go-seo/pkg/logger"

	"github.com/gin-gonic/gin"
)

type TrackingJobHandler struct {
	trackingJobUseCase *usecases.TrackingJobUseCase
}

func NewTrackingJobHandler(trackingJobUseCase *usecases.TrackingJobUseCase) *TrackingJobHandler {
	return &TrackingJobHandler{
		trackingJobUseCase: trackingJobUseCase,
	}
}

// GetTrackingJobs godoc
// @Summary Получить список джобов с пагинацией
// @Description Возвращает постраничный список джобов отслеживания позиций с возможностью фильтрации по сайту и статусу
// @Tags tracking-jobs
// @Accept json
// @Produce json
// @Param site_id query int false "ID сайта для фильтрации"
// @Param status query string false "Статус джоба (pending, running, completed, failed, cancelled)"
// @Param page query int false "Номер страницы (по умолчанию 1)"
// @Param per_page query int false "Количество записей на странице (по умолчанию 20, максимум 100)"
// @Success 200 {object} dto.TrackingJobsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/tracking-jobs [get]
func (h *TrackingJobHandler) GetTrackingJobs(c *gin.Context) {
	startTime := time.Now()

	var req dto.TrackingJobsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		logger.ErrorLogger.Printf("Failed to bind query parameters: %v", err)
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Получаем данные из use case
	response, err := h.trackingJobUseCase.GetJobsWithPagination(&req)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to get tracking jobs: %v", err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to retrieve tracking jobs",
		})
		return
	}

	// Обновляем время выполнения запроса
	response.Meta.QueryTimeMs = int(time.Since(startTime).Milliseconds())

	logger.InfoLogger.Printf("Tracking jobs retrieved successfully - Total: %d, Page: %d, PerPage: %d, QueryTime: %dms",
		response.Pagination.Total, response.Pagination.CurrentPage, response.Pagination.PerPage, response.Meta.QueryTimeMs)

	c.JSON(http.StatusOK, response)
}
