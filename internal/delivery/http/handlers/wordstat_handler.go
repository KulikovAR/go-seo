package handlers

import (
	"net/http"
	"strconv"

	"go-seo/internal/delivery/http/dto"
	"go-seo/internal/usecases"

	"github.com/gin-gonic/gin"
)

type WordstatHandler struct {
	wordstatUseCase *usecases.WordstatUseCase
}

func NewWordstatHandler(wordstatUseCase *usecases.WordstatUseCase) *WordstatHandler {
	return &WordstatHandler{
		wordstatUseCase: wordstatUseCase,
	}
}

// TrackKeywordFrequency godoc
// @Summary Track keyword frequency from Wordstat
// @Description Track frequency for specific keyword from Yandex Wordstat
// @Tags wordstat
// @Accept json
// @Produce json
// @Param keyword_id path int true "Keyword ID"
// @Success 200 {object} dto.WordstatFrequencyResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wordstat/keyword/{keyword_id}/frequency [post]
func (h *WordstatHandler) TrackKeywordFrequency(c *gin.Context) {
	keywordIDStr := c.Param("keyword_id")
	keywordID, err := strconv.Atoi(keywordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid keyword_id parameter",
		})
		return
	}

	frequency, err := h.wordstatUseCase.TrackKeywordFrequency(keywordID)
	if err != nil {
		if usecases.IsDomainError(err) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   usecases.GetDomainErrorCode(err),
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to track keyword frequency",
		})
		return
	}

	c.JSON(http.StatusOK, dto.WordstatFrequencyResponse{
		Message:   "Keyword frequency tracked successfully",
		Frequency: frequency,
	})
}

// TrackSiteKeywordsFrequency godoc
// @Summary Track keywords frequency for site from Wordstat
// @Description Track frequency for all keywords of specific site from Yandex Wordstat
// @Tags wordstat
// @Accept json
// @Produce json
// @Param site_id path int true "Site ID"
// @Success 200 {object} dto.TrackPositionsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wordstat/site/{site_id}/frequency [post]
func (h *WordstatHandler) TrackSiteKeywordsFrequency(c *gin.Context) {
	siteIDStr := c.Param("site_id")
	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid site_id parameter",
		})
		return
	}

	count, err := h.wordstatUseCase.TrackSiteKeywordsFrequency(siteID)
	if err != nil {
		if usecases.IsDomainError(err) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   usecases.GetDomainErrorCode(err),
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to track site keywords frequency",
		})
		return
	}

	c.JSON(http.StatusOK, dto.TrackPositionsResponse{
		Message: "Site keywords frequency tracked successfully",
		Count:   count,
	})
}

// GetRelatedKeywords godoc
// @Summary Get related keywords from Wordstat
// @Description Get related keywords and associations from Yandex Wordstat
// @Tags wordstat
// @Accept json
// @Produce json
// @Param query query string true "Search query"
// @Success 200 {object} dto.WordstatRelatedResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wordstat/related [get]
func (h *WordstatHandler) GetRelatedKeywords(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "query parameter is required",
		})
		return
	}

	relatedKeywords, err := h.wordstatUseCase.GetRelatedKeywords(query)
	if err != nil {
		if usecases.IsDomainError(err) {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   usecases.GetDomainErrorCode(err),
				Message: err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to get related keywords",
		})
		return
	}

	var response []dto.WordstatItemResponse
	for _, item := range relatedKeywords {
		response = append(response, dto.WordstatItemResponse{
			IsAssociations: item.IsAssociations,
			Value:          item.Value,
			Text:           item.Text,
		})
	}

	c.JSON(http.StatusOK, dto.WordstatRelatedResponse{
		Message: "Related keywords retrieved successfully",
		Items:   response,
	})
}
