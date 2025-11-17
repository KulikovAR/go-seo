package usecases

import (
	"go-seo/internal/infrastructure/services"
	"go-seo/internal/repositories"
)

type Container struct {
	Site                  *SiteUseCase
	Keyword               *KeywordUseCase
	Group                 *GroupUseCase
	PositionTracking      *PositionTrackingUseCase
	AsyncPositionTracking *AsyncPositionTrackingUseCase
	TrackingJob           *TrackingJobUseCase
	Debug                 *DebugUseCase
}

func NewContainer(repos *repositories.Container, xmlRiver *services.XMLRiverService, xmlStock *services.XMLRiverService, wordstat *services.WordstatService, kafkaService *services.KafkaService, idGenerator *services.IDGeneratorService, retryService *services.RetryService, workerCount int, batchSize int, xmlRiverSoftID string, xmlStockSoftID string) *Container {
	return &Container{
		Site:                  NewSiteUseCase(repos.Site, repos.Position, repos.Keyword, repos.Group, repos.TrackingJob, repos.TrackingTask, repos.TrackingResult),
		Keyword:               NewKeywordUseCase(repos.Keyword, repos.Position),
		Group:                 NewGroupUseCase(repos.Group),
		PositionTracking:      NewPositionTrackingUseCase(repos.Site, repos.Keyword, repos.Position, xmlRiver, xmlStock, wordstat, xmlRiverSoftID, xmlStockSoftID),
		AsyncPositionTracking: NewAsyncPositionTrackingUseCase(repos.Site, repos.Keyword, repos.Position, repos.TrackingJob, repos.TrackingTask, repos.TrackingResult, xmlRiver, xmlStock, wordstat, kafkaService, idGenerator, retryService, workerCount, batchSize, xmlRiverSoftID, xmlStockSoftID),
		TrackingJob:           NewTrackingJobUseCase(repos.TrackingJob),
		Debug:                 NewDebugUseCase(kafkaService),
	}
}
