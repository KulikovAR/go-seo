package repositories

import (
	"go-seo/internal/domain/repositories"

	"gorm.io/gorm"
)

type RepositoryContainer struct {
	Keyword        repositories.KeywordRepository
	Site           repositories.SiteRepository
	Group          repositories.GroupRepository
	Position       repositories.PositionRepository
	TrackingJob    repositories.TrackingJobRepository
	TrackingTask   repositories.TrackingTaskRepository
	TrackingResult repositories.TrackingResultRepository
}

func NewRepositoryContainer(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		Keyword:        NewKeywordRepository(db),
		Site:           NewSiteRepository(db),
		Group:          NewGroupRepository(db),
		Position:       NewPositionRepository(db),
		TrackingJob:    NewTrackingJobRepository(db),
		TrackingTask:   NewTrackingTaskRepository(db),
		TrackingResult: NewTrackingResultRepository(db),
	}
}
