package handlers

import (
	"net/http"

	"go-seo/internal/usecases"

	"github.com/gin-gonic/gin"
)

type DebugHandler struct {
	debugUseCase *usecases.DebugUseCase
}

func NewDebugHandler(debugUseCase *usecases.DebugUseCase) *DebugHandler {
	return &DebugHandler{
		debugUseCase: debugUseCase,
	}
}

type kafkaJobStatusRequest struct {
	JobID   string `json:"job_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
	Error   string `json:"error"`
	Percent *int   `json:"percent"`
}

func (h *DebugHandler) SendKafkaJobStatus(c *gin.Context) {
	var req kafkaJobStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := h.debugUseCase.SendJobStatus(req.JobID, req.Status, req.Error, req.Percent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kafka job status message sent",
	})
}

