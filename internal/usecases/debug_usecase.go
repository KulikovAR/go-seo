package usecases

import (
	"fmt"

	"go-seo/internal/infrastructure/services"
)

type DebugUseCase struct {
	kafkaService *services.KafkaService
}

func NewDebugUseCase(kafkaService *services.KafkaService) *DebugUseCase {
	return &DebugUseCase{
		kafkaService: kafkaService,
	}
}

func (uc *DebugUseCase) SendJobStatus(jobID, status, errorMsg string, percent *int) error {
	if jobID == "" {
		return fmt.Errorf("job_id is required")
	}
	if status == "" {
		return fmt.Errorf("status is required")
	}

	if percent != nil {
		return uc.kafkaService.SendJobStatus(jobID, status, errorMsg, *percent)
	}

	return uc.kafkaService.SendJobStatus(jobID, status, errorMsg)
}
