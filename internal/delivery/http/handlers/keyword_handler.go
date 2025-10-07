package handlers

import (
	"net/http"
	"strconv"

	"go-seo/internal/delivery/http/dto"
	"go-seo/internal/usecases"

	"github.com/gin-gonic/gin"
)

type KeywordHandler struct {
	keywordUseCase usecases.KeywordUseCaseInterface
}

func NewKeywordHandler(keywordUseCase usecases.KeywordUseCaseInterface) *KeywordHandler {
	return &KeywordHandler{
		keywordUseCase: keywordUseCase,
	}
}

// CreateKeyword godoc
// @Summary Create a new keyword
// @Description Create a new keyword for tracking
// @Tags keywords
// @Accept json
// @Produce json
// @Param keyword body dto.CreateKeywordRequest true "Keyword data"
// @Success 201 {object} dto.KeywordResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/keywords [post]
func (h *KeywordHandler) CreateKeyword(c *gin.Context) {
	var req dto.CreateKeywordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	keyword, err := h.keywordUseCase.CreateKeyword(req.Value, req.SiteID)
	if err != nil {
		if usecases.IsDomainError(err) {
			code := usecases.GetDomainErrorCode(err)
			status := http.StatusInternalServerError

			switch code {
			case usecases.ErrorKeywordExists:
				status = http.StatusConflict
			case usecases.ErrorKeywordCreation:
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

	c.JSON(http.StatusCreated, dto.KeywordResponse{
		ID:     keyword.ID,
		Value:  keyword.Value,
		SiteID: keyword.SiteID,
	})
}

// DeleteKeyword godoc
// @Summary Delete a keyword
// @Description Delete a keyword and all its tracking data
// @Tags keywords
// @Param id path int true "Keyword ID"
// @Success 200 {object} dto.DeleteKeywordResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/keywords/{id} [delete]
func (h *KeywordHandler) DeleteKeyword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid keyword ID",
		})
		return
	}

	err = h.keywordUseCase.DeleteKeyword(id)
	if err != nil {
		if usecases.IsDomainError(err) {
			code := usecases.GetDomainErrorCode(err)
			status := http.StatusInternalServerError

			switch code {
			case usecases.ErrorKeywordNotFound:
				status = http.StatusNotFound
			case usecases.ErrorKeywordDeletion, usecases.ErrorPositionDeletion:
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

	c.JSON(http.StatusOK, dto.DeleteKeywordResponse{
		Message: "Keyword and all tracking data deleted successfully",
	})
}

// GetKeywords godoc
// @Summary Get keywords by site
// @Description Get list of keywords for a specific site
// @Tags keywords
// @Produce json
// @Param site_id query int true "Site ID"
// @Success 200 {array} dto.KeywordResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/keywords [get]
func (h *KeywordHandler) GetKeywords(c *gin.Context) {
	siteIDStr := c.Query("site_id")
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

	keywords, err := h.keywordUseCase.GetKeywordsBySite(siteID)
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

	response := make([]dto.KeywordResponse, len(keywords))
	for i, keyword := range keywords {
		response[i] = dto.KeywordResponse{
			ID:     keyword.ID,
			Value:  keyword.Value,
			SiteID: keyword.SiteID,
		}
	}

	c.JSON(http.StatusOK, response)
}
