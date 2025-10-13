package handlers

import (
	"net/http"
	"strconv"

	"go-seo/internal/delivery/http/dto"
	"go-seo/internal/usecases"

	"github.com/gin-gonic/gin"
)

type PositionHandler struct {
	positionTrackingUseCase *usecases.PositionTrackingUseCase
}

func NewPositionHandler(positionTrackingUseCase *usecases.PositionTrackingUseCase) *PositionHandler {
	return &PositionHandler{
		positionTrackingUseCase: positionTrackingUseCase,
	}
}

func (h *PositionHandler) TrackSitePositions(c *gin.Context) {
	var req dto.TrackSitePositionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if req.Device == "mobile" && req.OS == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "OS parameter is required when device is mobile",
		})
		return
	}

	count, err := h.positionTrackingUseCase.TrackSitePositions(
		req.SiteID,
		req.Source,
		req.Device,
		req.OS,
		req.Ads,
		req.Country,
		req.Lang,
		req.Pages,
	)

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
			Message: "Failed to track site positions",
		})
		return
	}

	c.JSON(http.StatusOK, dto.TrackPositionsResponse{
		Message: "Site positions tracked successfully",
		Count:   count,
	})
}

func (h *PositionHandler) GetPositionsHistory(c *gin.Context) {
	siteIDStr := c.Query("site_id")
	keywordIDStr := c.Query("keyword_id")

	if siteIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "site_id parameter is required",
		})
		return
	}

	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid site_id parameter",
		})
		return
	}

	var keywordID *int
	if keywordIDStr != "" {
		id, err := strconv.Atoi(keywordIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid keyword_id parameter",
			})
			return
		}
		keywordID = &id
	}

	positions, err := h.positionTrackingUseCase.GetPositionsHistory(siteID, keywordID)
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
			Message: "Failed to fetch positions history",
		})
		return
	}

	var response []dto.PositionResponse
	for _, pos := range positions {
		keywordValue := ""
		if pos.Keyword != nil {
			keywordValue = pos.Keyword.Value
		}
		response = append(response, dto.PositionResponse{
			ID:        pos.ID,
			KeywordID: pos.KeywordID,
			SiteID:    pos.SiteID,
			Rank:      pos.Rank,
			URL:       pos.URL,
			Title:     pos.Title,
			Source:    pos.Source,
			Device:    pos.Device,
			OS:        pos.OS,
			Ads:       pos.Ads,
			Country:   pos.Country,
			Lang:      pos.Lang,
			Pages:     pos.Pages,
			Date:      pos.Date,
			Keyword:   keywordValue,
		})
	}

	c.JSON(http.StatusOK, response)
}

func (h *PositionHandler) GetLatestPositions(c *gin.Context) {
	positions, err := h.positionTrackingUseCase.GetLatestPositions()
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
			Message: "Failed to fetch latest positions",
		})
		return
	}

	var response []dto.PositionResponse
	for _, pos := range positions {
		response = append(response, dto.PositionResponse{
			ID:        pos.ID,
			KeywordID: pos.KeywordID,
			SiteID:    pos.SiteID,
			Rank:      pos.Rank,
			URL:       pos.URL,
			Title:     pos.Title,
			Source:    pos.Source,
			Device:    pos.Device,
			OS:        pos.OS,
			Ads:       pos.Ads,
			Country:   pos.Country,
			Lang:      pos.Lang,
			Pages:     pos.Pages,
			Date:      pos.Date,
		})
	}

	c.JSON(http.StatusOK, response)
}
