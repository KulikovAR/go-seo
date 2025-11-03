package handlers

import (
	"net/http"
	"time"

	"go-seo/internal/delivery/http/dto"
	"go-seo/internal/usecases"
	"go-seo/pkg/logger"

	"github.com/gin-gonic/gin"
)

type PositionHandler struct {
	positionTrackingUseCase      *usecases.PositionTrackingUseCase
	asyncPositionTrackingUseCase *usecases.AsyncPositionTrackingUseCase
}

func NewPositionHandler(positionTrackingUseCase *usecases.PositionTrackingUseCase, asyncPositionTrackingUseCase *usecases.AsyncPositionTrackingUseCase) *PositionHandler {
	return &PositionHandler{
		positionTrackingUseCase:      positionTrackingUseCase,
		asyncPositionTrackingUseCase: asyncPositionTrackingUseCase,
	}
}

// @Summary Track Google positions
// @Description Start async Google position tracking for site keywords
// @Accept json
// @Produce json
// @Param request body dto.TrackGooglePositionsRequest true "Google tracking parameters"
// @Success 200 {object} dto.AsyncTrackPositionsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/positions/track-google [post]
func (h *PositionHandler) TrackGooglePositions(c *gin.Context) {
	var req dto.TrackGooglePositionsRequest
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

	logger.LogTrackSiteParams(
		req.SiteID,
		"google",
		req.Device,
		req.OS,
		req.Ads,
		req.Country,
		req.Lang,
		req.Pages,
		req.Subdomains,
		0,
	)

	taskID, err := h.asyncPositionTrackingUseCase.StartAsyncGoogleTracking(
		req.SiteID,
		req.Device,
		req.OS,
		req.Ads,
		req.Country,
		req.Lang,
		req.Pages,
		req.Subdomains,
		req.XMLUserID,
		req.XMLAPIKey,
		req.XMLBaseURL,
		req.TBS,
		req.Filter,
		req.Highlights,
		req.NFPR,
		req.Loc,
		req.AI,
		req.Raw,
		req.LR,
		req.Domain,
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
			Message: "Failed to start Google tracking",
		})
		return
	}

	c.JSON(http.StatusOK, dto.AsyncTrackPositionsResponse{
		Message: "Google tracking started successfully",
		TaskID:  taskID,
		Status:  "pending",
	})
}

// @Summary Track Yandex positions
// @Description Start async Yandex position tracking for site keywords
// @Accept json
// @Produce json
// @Param request body dto.TrackYandexPositionsRequest true "Yandex tracking parameters"
// @Success 200 {object} dto.AsyncTrackPositionsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/positions/track-yandex [post]
func (h *PositionHandler) TrackYandexPositions(c *gin.Context) {
	var req dto.TrackYandexPositionsRequest
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

	logger.LogTrackSiteParams(
		req.SiteID,
		"yandex",
		req.Device,
		req.OS,
		req.Ads,
		req.Country,
		req.Lang,
		req.Pages,
		req.Subdomains,
		req.LR,
	)

	taskID, err := h.asyncPositionTrackingUseCase.StartAsyncYandexTracking(
		req.SiteID,
		req.Device,
		req.OS,
		req.Ads,
		req.Country,
		req.Lang,
		req.Pages,
		req.Subdomains,
		req.XMLUserID,
		req.XMLAPIKey,
		req.XMLBaseURL,
		req.GroupBy,
		req.Filter,
		req.Highlights,
		req.Within,
		req.LR,
		req.Raw,
		req.InIndex,
		req.Strict,
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
			Message: "Failed to start Yandex tracking",
		})
		return
	}

	c.JSON(http.StatusOK, dto.AsyncTrackPositionsResponse{
		Message: "Yandex tracking started successfully",
		TaskID:  taskID,
		Status:  "pending",
	})
}

// @Summary Track Wordstat positions
// @Description Start async Wordstat position tracking for site keywords
// @Accept json
// @Produce json
// @Param request body dto.TrackWordstatPositionsRequest true "Wordstat tracking parameters"
// @Success 200 {object} dto.AsyncTrackPositionsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/positions/track-wordstat [post]
func (h *PositionHandler) TrackWordstatPositions(c *gin.Context) {
	var req dto.TrackWordstatPositionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	logger.LogTrackSiteParamsWithRegions(
		req.SiteID,
		"wordstat",
		"",
		"",
		false,
		"",
		"",
		0,
		false,
		req.Regions,
	)

	taskID, err := h.asyncPositionTrackingUseCase.StartAsyncWordstatTracking(
		req.SiteID,
		req.XMLUserID,
		req.XMLAPIKey,
		req.XMLBaseURL,
		req.Regions,
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
			Message: "Failed to start Wordstat tracking",
		})
		return
	}

	c.JSON(http.StatusOK, dto.AsyncTrackPositionsResponse{
		Message: "Wordstat tracking started successfully",
		TaskID:  taskID,
		Status:  "pending",
	})
}

// @Summary Get positions history
// @Description Get paginated positions history with filtering options
// @Accept json
// @Produce json
// @Param site_id query int true "Site ID"
// @Param keyword_id query int false "Keyword ID"
// @Param source query string false "Source (google, yandex, wordstat)"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param last query bool false "Get only last positions"
// @Param page query int false "Page number (default 1)"
// @Param per_page query int false "Items per page (default 50, max 100)"
// @Success 200 {object} dto.PositionHistoryResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/positions/history [get]
func (h *PositionHandler) GetPositionsHistory(c *gin.Context) {
	startTime := time.Now()

	var req dto.PositionHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 50
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}

	if req.Source != nil {
		if *req.Source != "google" && *req.Source != "yandex" && *req.Source != "wordstat" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "source must be either 'google', 'yandex' or 'wordstat'",
			})
			return
		}
	}

	var dateFrom, dateTo *time.Time
	if req.DateFrom != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid date_from parameter. Use YYYY-MM-DD format",
			})
			return
		}
		dateFrom = &parsed
	}
	if req.DateTo != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid date_to parameter. Use YYYY-MM-DD format",
			})
			return
		}
		dateTo = &parsed
	}

	var last bool
	if req.Last != nil {
		last = *req.Last
	}

	positions, total, err := h.positionTrackingUseCase.GetPositionsHistoryPaginated(
		req.SiteID, req.KeywordID, req.Source, dateFrom, dateTo, last, req.Page, req.PerPage)
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

	lastPage := int((total + int64(req.PerPage) - 1) / int64(req.PerPage))
	from := (req.Page-1)*req.PerPage + 1
	to := from + len(positions) - 1
	if len(positions) == 0 {
		from = 0
		to = 0
	}
	hasMore := req.Page < lastPage

	var data []dto.PositionHistoryItem
	for _, pos := range positions {
		keywordValue := ""
		if pos.Keyword != nil {
			keywordValue = pos.Keyword.Value
		}
		data = append(data, dto.PositionHistoryItem{
			ID:        pos.ID,
			SiteID:    pos.SiteID,
			KeywordID: pos.KeywordID,
			Keyword:   keywordValue,
			Rank:      pos.Rank,
			URL:       pos.URL,
			Title:     pos.Title,
			Date:      pos.Date,
			Source:    pos.Source,
			Device:    pos.Device,
			Country:   pos.Country,
			Lang:      pos.Lang,
		})
	}

	queryTimeMs := int(time.Since(startTime).Milliseconds())

	response := dto.PositionHistoryResponse{
		Data: data,
		Pagination: dto.PaginationInfo{
			CurrentPage: req.Page,
			PerPage:     req.PerPage,
			Total:       int(total),
			LastPage:    lastPage,
			From:        from,
			To:          to,
			HasMore:     hasMore,
		},
		Meta: dto.MetaInfo{
			QueryTimeMs: queryTimeMs,
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get combined positions
// @Description Get paginated combined positions from multiple sources
// @Accept json
// @Produce json
// @Param site_id query int true "Site ID"
// @Param source query string false "Source (google, yandex)"
// @Param wordstat query bool false "Include Wordstat data"
// @Param wordstat_sort query string false "Sort by Wordstat positions (asc or desc)"
// @Param date_from query string false "Start date (YYYY-MM-DD)"
// @Param date_to query string false "End date (YYYY-MM-DD)"
// @Param date_sort query string false "Date for sorting by position (YYYY-MM-DD). Must be within date_from and date_to range"
// @Param sort_type query string false "Sort type for positions (asc or desc). Default: asc"
// @Param rank_from query int false "Minimum rank filter"
// @Param rank_to query int false "Maximum rank filter"
// @Param page query int false "Page number (default 1)"
// @Param per_page query int false "Items per page (default 50, max 100)"
// @Success 200 {object} dto.CombinedPositionsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/positions/combined [get]
func (h *PositionHandler) GetCombinedPositions(c *gin.Context) {
	startTime := time.Now()

	var req dto.CombinedPositionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 50
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}

	if req.Source != nil {
		if *req.Source != "google" && *req.Source != "yandex" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "source must be either 'google' or 'yandex'",
			})
			return
		}
	}

	var dateFrom, dateTo *time.Time
	if req.DateFrom != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateFrom)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid date_from parameter. Use YYYY-MM-DD format",
			})
			return
		}
		dateFrom = &parsed
	}
	if req.DateTo != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateTo)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid date_to parameter. Use YYYY-MM-DD format",
			})
			return
		}
		dateTo = &parsed
	}

	var dateSort *time.Time
	if req.DateSort != nil {
		parsed, err := time.Parse("2006-01-02", *req.DateSort)
		if err != nil {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "Invalid date_sort parameter. Use YYYY-MM-DD format",
			})
			return
		}
		dateSort = &parsed

		if dateFrom != nil && dateTo != nil {
			if dateSort.Before(*dateFrom) {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "validation_error",
					Message: "date_sort must be greater than or equal to date_from",
				})
				return
			}
			if dateSort.After(*dateTo) {
				c.JSON(http.StatusBadRequest, dto.ErrorResponse{
					Error:   "validation_error",
					Message: "date_sort must be less than or equal to date_to",
				})
				return
			}
		}
	}

	sortType := "asc"
	if req.SortType != nil {
		if *req.SortType != "asc" && *req.SortType != "desc" {
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error:   "validation_error",
				Message: "sort_type must be either 'asc' or 'desc'",
			})
			return
		}
		sortType = *req.SortType
	}

	var includeWordstat bool
	if req.Wordstat != nil {
		includeWordstat = *req.Wordstat
	}

	var wordstatSort bool
	if req.WordstatSort != nil {
		wordstatSort = true
		sortType = *req.WordstatSort
	}

	combinedPositions, total, err := h.positionTrackingUseCase.GetCombinedPositionsPaginated(
		req.SiteID, req.Source, includeWordstat, wordstatSort, dateFrom, dateTo, dateSort, sortType, req.RankFrom, req.RankTo, req.Page, req.PerPage)
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
			Message: "Failed to fetch combined positions",
		})
		return
	}

	lastPage := int((total + int64(req.PerPage) - 1) / int64(req.PerPage))
	from := (req.Page-1)*req.PerPage + 1
	to := from + len(combinedPositions) - 1
	if len(combinedPositions) == 0 {
		from = 0
		to = 0
	}
	hasMore := req.Page < lastPage

	var data []dto.CombinedPositionItem
	for _, pos := range combinedPositions {
		keywordValue := ""
		if pos.Keyword != nil {
			keywordValue = pos.Keyword.Value
		}

		item := dto.CombinedPositionItem{
			ID:        pos.ID,
			SiteID:    pos.SiteID,
			KeywordID: pos.KeywordID,
			Keyword:   keywordValue,
			Date:      pos.Date,
		}

		for _, position := range pos.Positions {
			item.Positions = append(item.Positions, dto.PositionData{
				Rank:   position.Rank,
				URL:    position.URL,
				Title:  position.Title,
				Source: position.Source,
				Date:   position.Date,
			})
		}

		if pos.Wordstat != nil {
			item.Wordstat = &dto.PositionData{
				Rank:   pos.Wordstat.Rank,
				URL:    pos.Wordstat.URL,
				Title:  pos.Wordstat.Title,
				Source: pos.Wordstat.Source,
				Date:   pos.Wordstat.Date,
			}
		}

		data = append(data, item)
	}

	queryTimeMs := int(time.Since(startTime).Milliseconds())

	response := dto.CombinedPositionsResponse{
		Data: data,
		Pagination: dto.PaginationInfo{
			CurrentPage: req.Page,
			PerPage:     req.PerPage,
			Total:       int(total),
			LastPage:    lastPage,
			From:        from,
			To:          to,
			HasMore:     hasMore,
		},
		Meta: dto.MetaInfo{
			QueryTimeMs: queryTimeMs,
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get position statistics
// @Description Get position statistics for a site within date range
// @Accept json
// @Produce json
// @Param request body dto.PositionStatisticsRequest true "Statistics parameters"
// @Success 200 {object} dto.PositionStatisticsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/positions/statistics [post]
func (h *PositionHandler) GetPositionStatistics(c *gin.Context) {
	var req dto.PositionStatisticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	dateFrom, err := time.Parse("2006-01-02", req.DateFrom)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid date_from format. Use YYYY-MM-DD",
		})
		return
	}

	dateTo, err := time.Parse("2006-01-02", req.DateTo)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid date_to format. Use YYYY-MM-DD",
		})
		return
	}

	stats, err := h.positionTrackingUseCase.GetPositionStatistics(req.SiteID, req.Source, dateFrom, dateTo)
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
			Message: "Failed to fetch position statistics",
		})
		return
	}

	response := dto.PositionStatisticsResponse{
		TotalPositions: stats.TotalPositions,
		KeywordsCount:  stats.KeywordsCount,
		Visible:        stats.Visible,
		NotVisible:     stats.NotVisible,
		PositionDistribution: dto.PositionDistribution{
			Top3:     stats.PositionDistribution.Top3,
			Top10:    stats.PositionDistribution.Top10,
			Top20:    stats.PositionDistribution.Top20,
			NotFound: stats.PositionDistribution.NotFound,
		},
		PositionRanges: dto.PositionRanges{
			Range1_3:     stats.PositionRanges.Range1_3,
			Range4_10:    stats.PositionRanges.Range4_10,
			Range11_30:   stats.PositionRanges.Range11_30,
			Range31_50:   stats.PositionRanges.Range31_50,
			Range51_100:  stats.PositionRanges.Range51_100,
			Range100Plus: stats.PositionRanges.Range100Plus,
			NotFound:     stats.PositionRanges.NotFound,
		},
		VisibilityStats: dto.VisibilityStats{
			AvgPosition:    stats.VisibilityStats.AvgPosition,
			MedianPosition: stats.VisibilityStats.MedianPosition,
			BestPosition:   stats.VisibilityStats.BestPosition,
			WorstPosition:  stats.VisibilityStats.WorstPosition,
		},
		Trends: dto.Trends{
			Improved: stats.Trends.Improved,
			Declined: stats.Trends.Declined,
			Stable:   stats.Trends.Stable,
		},
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get latest positions
// @Description Get latest positions for all keywords
// @Accept json
// @Produce json
// @Success 200 {array} dto.PositionResponse
// @Failure 400 {object} dto.ErrorResponse
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
