package handlers

import (
	"net/http"
	"strconv"
	"time"

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

// TrackSitePositions godoc
// @Summary Track positions for specific site
// @Description Track positions for specific site and its keywords. Supports Google, Yandex and Wordstat sources. Can include subdomains in search.
// @Tags positions
// @Accept json
// @Produce json
// @Param request body dto.TrackSitePositionsRequest true "Site tracking parameters (source: google, yandex or wordstat)"
// @Success 200 {object} dto.TrackPositionsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/positions/track-site [post]
func (h *PositionHandler) TrackSitePositions(c *gin.Context) {
	var req dto.TrackSitePositionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if req.Source != "wordstat" {
		if req.Pages == 0 {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "pages parameter is required for google and yandex sources",
			})
			return
		}
		if req.Device == "" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "device parameter is required for google and yandex sources",
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
	}

	device := req.Device
	pages := req.Pages
	if req.Source == "wordstat" {
		device = ""
		pages = 0
	}

	count, err := h.positionTrackingUseCase.TrackSitePositions(
		req.SiteID,
		req.Source,
		device,
		req.OS,
		req.Ads,
		req.Country,
		req.Lang,
		pages,
		req.Subdomains,
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

// GetPositionsHistory godoc
// @Summary Get positions history
// @Description Get positions history for specific site and optional keyword
// @Tags positions
// @Produce json
// @Param site_id query int true "Site ID"
// @Param keyword_id query int false "Keyword ID (optional)"
// @Param source query string false "Source filter (optional) - google, yandex or wordstat" Enums(google,yandex,wordstat)
// @Param date_from query string false "Start date filter (optional) - YYYY-MM-DD format"
// @Param date_to query string false "End date filter (optional) - YYYY-MM-DD format"
// @Success 200 {array} dto.PositionResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/positions/history [get]
func (h *PositionHandler) GetPositionsHistory(c *gin.Context) {
	siteIDStr := c.Query("site_id")
	keywordIDStr := c.Query("keyword_id")
	sourceStr := c.Query("source")
	dateFromStr := c.Query("date_from")
	dateToStr := c.Query("date_to")

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

	var source *string
	if sourceStr != "" {
		if sourceStr != "google" && sourceStr != "yandex" && sourceStr != "wordstat" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "source must be either 'google', 'yandex' or 'wordstat'",
			})
			return
		}
		source = &sourceStr
	}

	var dateFrom, dateTo *time.Time
	if dateFromStr != "" {
		parsed, err := time.Parse("2006-01-02", dateFromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid date_from parameter. Use YYYY-MM-DD format",
			})
			return
		}
		dateFrom = &parsed
	}
	if dateToStr != "" {
		parsed, err := time.Parse("2006-01-02", dateToStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid date_to parameter. Use YYYY-MM-DD format",
			})
			return
		}
		dateTo = &parsed
	}

	positions, err := h.positionTrackingUseCase.GetPositionsHistory(siteID, keywordID, source, dateFrom, dateTo)
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

// GetLatestPositions godoc
// @Summary Get latest positions
// @Description Get latest positions for all sites and keywords
// @Tags positions
// @Produce json
// @Success 200 {array} dto.PositionResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/positions/latest [get]
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
