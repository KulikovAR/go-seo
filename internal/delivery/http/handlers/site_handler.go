package handlers

import (
	"net/http"
	"strconv"

	"go-seo/internal/delivery/http/dto"
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
		ID:     site.ID,
		Domain: site.Domain,
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
// @Summary Get all sites
// @Description Get list of all tracked sites
// @Tags sites
// @Produce json
// @Success 200 {array} dto.SiteResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/sites [get]
func (h *SiteHandler) GetSites(c *gin.Context) {
	sites, err := h.siteUseCase.GetAllSites()
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
		response[i] = dto.SiteResponse{
			ID:     site.ID,
			Domain: site.Domain,
		}
	}

	c.JSON(http.StatusOK, response)
}
