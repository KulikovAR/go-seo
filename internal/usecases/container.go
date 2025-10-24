package usecases

import (
	"go-seo/internal/infrastructure/services"
	"go-seo/internal/repositories"
)

type Container struct {
	Site                  *SiteUseCase
	Keyword               *KeywordUseCase
	PositionTracking      *PositionTrackingUseCase
	AsyncPositionTracking *AsyncPositionTrackingUseCase
	TrackingJob           *TrackingJobUseCase
}

func NewContainer(repos *repositories.Container, xmlRiver *services.XMLRiverService, xmlStock *services.XMLRiverService, wordstat *services.WordstatService, kafkaService *services.KafkaService, idGenerator *services.IDGeneratorService, retryService *services.RetryService, workerCount int, batchSize int) *Container {
	return &Container{
		Site:                  NewSiteUseCase(repos.Site, repos.Position, repos.Keyword),
		Keyword:               NewKeywordUseCase(repos.Keyword, repos.Position),
		PositionTracking:      NewPositionTrackingUseCase(repos.Site, repos.Keyword, repos.Position, xmlRiver, xmlStock, wordstat),
		AsyncPositionTracking: NewAsyncPositionTrackingUseCase(repos.Site, repos.Keyword, repos.Position, repos.TrackingJob, repos.TrackingTask, repos.TrackingResult, xmlRiver, xmlStock, wordstat, kafkaService, idGenerator, retryService, workerCount, batchSize),
		TrackingJob:           NewTrackingJobUseCase(repos.TrackingJob),
	}
}
