package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"go-seo/internal/delivery/http/dto"
	"go-seo/internal/domain/entities"
	"go-seo/internal/usecases"

	"github.com/gin-gonic/gin"
)

type SiteHandler struct {
	siteUseCase usecases.SiteUseCaseInterface
}

func NewSiteHandler(siteUseCase usecases.SiteUseCaseInterface) *SiteHandler {
	return &SiteHandler{
		siteUseCase: siteUseCase,
	}
}

func parseIDsFromQuery(c *gin.Context) ([]int, error) {
	idsStr := c.Query("ids")
	if idsStr == "" {
		return nil, nil
	}

	idStrings := strings.Split(idsStr, ",")
	ids := make([]int, 0, len(idStrings))

	for _, idStr := range idStrings {
		idStr = strings.TrimSpace(idStr)
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// CreateSite godoc
// @Summary Create a new site
// @Description Create a new site for tracking
// @Tags sites
// @Accept json
// @Produce json
// @Param site body dto.CreateSiteRequest true "Site data"
// @Success 201 {object} dto.SiteResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/sites [post]
func (h *SiteHandler) CreateSite(c *gin.Context) {
	var req dto.CreateSiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	site, err := h.siteUseCase.CreateSite(req.Domain)
	if err != nil {
		if usecases.IsDomainError(err) {
			code := usecases.GetDomainErrorCode(err)
			status := http.StatusInternalServerError

			switch code {
			case usecases.ErrorSiteExists:
				status = http.StatusConflict
			case usecases.ErrorSiteCreation:
				status = http.StatusInternalServerError
			}

			c.JSON(status, dto.ErrorResponse{
				Error:   code,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SiteResponse{
		ID:                 site.ID,
		Domain:             site.Domain,
		KeywordsCount:      0,
		LastPositionUpdate: nil,
	})
}

// DeleteSite godoc
// @Summary Delete a site
// @Description Delete a site and all its tracking data
// @Tags sites
// @Param id path int true "Site ID"
// @Success 200 {object} dto.DeleteSiteResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/sites/{id} [delete]
func (h *SiteHandler) DeleteSite(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid site ID",
		})
		return
	}

	err = h.siteUseCase.DeleteSite(id)
	if err != nil {
		if usecases.IsDomainError(err) {
			code := usecases.GetDomainErrorCode(err)
			status := http.StatusInternalServerError

			switch code {
			case usecases.ErrorSiteNotFound:
				status = http.StatusNotFound
			case usecases.ErrorSiteDeletion, usecases.ErrorPositionDeletion:
				status = http.StatusInternalServerError
			}

			c.JSON(status, dto.ErrorResponse{
				Error:   code,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, dto.DeleteSiteResponse{
		Message: "Site and all tracking data deleted successfully",
	})
}

// GetSites godoc
// @Summary Get sites
// @Description Get list of tracked sites. If ids parameter is provided, returns only sites with specified IDs
// @Tags sites
// @Produce json
// @Param ids query string false "Comma-separated list of site IDs to filter by"
// @Success 200 {array} dto.SiteResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/sites [get]
func (h *SiteHandler) GetSites(c *gin.Context) {
	ids, err := parseIDsFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_ids",
			Message: "Invalid IDs format. Expected comma-separated integers",
		})
		return
	}

	var sites []*entities.Site
	if ids != nil {
		sites, err = h.siteUseCase.GetSitesByIDs(ids)
	} else {
		sites, err = h.siteUseCase.GetAllSites()
	}

	if err != nil {
		if usecases.IsDomainError(err) {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error:   usecases.GetDomainErrorCode(err),
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Internal server error",
		})
		return
	}

	response := make([]dto.SiteResponse, len(sites))
	for i, site := range sites {
		keywordsCount, err := h.siteUseCase.GetKeywordsCount(site.ID)
		if err != nil {
			keywordsCount = 0
		}

		lastPositionUpdate, err := h.siteUseCase.GetLastPositionUpdateDate(site.ID)
		if err != nil {
			lastPositionUpdate = nil
		}

		response[i] = dto.SiteResponse{
			ID:                 site.ID,
			Domain:             site.Domain,
			KeywordsCount:      keywordsCount,
			LastPositionUpdate: lastPositionUpdate,
		}
	}

	c.JSON(http.StatusOK, response)
}
