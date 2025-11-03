package handlers

import (
	"net/http"
	"strconv"

	"go-seo/internal/delivery/http/dto"
	"go-seo/internal/usecases"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	groupUseCase usecases.GroupUseCaseInterface
}

func NewGroupHandler(groupUseCase usecases.GroupUseCaseInterface) *GroupHandler {
	return &GroupHandler{
		groupUseCase: groupUseCase,
	}
}

// CreateGroup godoc
// @Summary Create a new group
// @Description Create a new group for organizing keywords
// @Tags groups
// @Accept json
// @Produce json
// @Param group body dto.CreateGroupRequest true "Group data"
// @Success 201 {object} dto.GroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/groups [post]
func (h *GroupHandler) CreateGroup(c *gin.Context) {
	var req dto.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	group, err := h.groupUseCase.CreateGroup(req.Name, req.SiteID)
	if err != nil {
		if usecases.IsDomainError(err) {
			code := usecases.GetDomainErrorCode(err)
			status := http.StatusInternalServerError

			switch code {
			case usecases.ErrorGroupExists:
				status = http.StatusConflict
			case usecases.ErrorGroupCreation:
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

	c.JSON(http.StatusCreated, dto.GroupResponse{
		ID:     group.ID,
		Name:   group.Name,
		SiteID: group.SiteID,
	})
}

// UpdateGroup godoc
// @Summary Update a group
// @Description Update group name
// @Tags groups
// @Accept json
// @Produce json
// @Param id path int true "Group ID"
// @Param group body dto.UpdateGroupRequest true "Group data"
// @Success 200 {object} dto.GroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/groups/{id} [put]
func (h *GroupHandler) UpdateGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid group ID",
		})
		return
	}

	var req dto.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	group, err := h.groupUseCase.UpdateGroup(id, req.Name)
	if err != nil {
		if usecases.IsDomainError(err) {
			code := usecases.GetDomainErrorCode(err)
			status := http.StatusInternalServerError

			switch code {
			case usecases.ErrorGroupNotFound:
				status = http.StatusNotFound
			case usecases.ErrorGroupFetch:
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

	c.JSON(http.StatusOK, dto.GroupResponse{
		ID:     group.ID,
		Name:   group.Name,
		SiteID: group.SiteID,
	})
}

// DeleteGroup godoc
// @Summary Delete a group
// @Description Delete a group
// @Tags groups
// @Param id path int true "Group ID"
// @Success 200 {object} dto.ErrorResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/groups/{id} [delete]
func (h *GroupHandler) DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_id",
			Message: "Invalid group ID",
		})
		return
	}

	err = h.groupUseCase.DeleteGroup(id)
	if err != nil {
		if usecases.IsDomainError(err) {
			code := usecases.GetDomainErrorCode(err)
			status := http.StatusInternalServerError

			switch code {
			case usecases.ErrorGroupNotFound:
				status = http.StatusNotFound
			case usecases.ErrorGroupDeletion:
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

	c.JSON(http.StatusOK, dto.ErrorResponse{
		Error:   "success",
		Message: "Group deleted successfully",
	})
}

// GetGroups godoc
// @Summary Get all groups
// @Description Get list of all groups for a specific site
// @Tags groups
// @Produce json
// @Param site_id query int true "Site ID"
// @Success 200 {array} dto.GroupResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/groups [get]
func (h *GroupHandler) GetGroups(c *gin.Context) {
	siteIDStr := c.Query("site_id")
	if siteIDStr == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "site_id is required",
		})
		return
	}

	siteID, err := strconv.Atoi(siteIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid site_id",
		})
		return
	}

	groups, err := h.groupUseCase.GetGroupsBySite(siteID)
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

	response := make([]dto.GroupResponse, len(groups))
	for i, group := range groups {
		response[i] = dto.GroupResponse{
			ID:     group.ID,
			Name:   group.Name,
			SiteID: group.SiteID,
		}
	}

	c.JSON(http.StatusOK, response)
}
