package repositories

import (
	"go-seo/internal/domain/repositories"
	postgresRepos "go-seo/internal/infrastructure/database/postgres/repositories"

	"gorm.io/gorm"
)

type Container struct {
	Keyword        repositories.KeywordRepository
	Site           repositories.SiteRepository
	Group          repositories.GroupRepository
	Position       repositories.PositionRepository
	TrackingJob    repositories.TrackingJobRepository
	TrackingTask   repositories.TrackingTaskRepository
	TrackingResult repositories.TrackingResultRepository
}

func NewContainer(db *gorm.DB) *Container {
	postgresRepos := postgresRepos.NewRepositoryContainer(db)

	return &Container{
		Keyword:        postgresRepos.Keyword,
		Site:           postgresRepos.Site,
		Group:          postgresRepos.Group,
		Position:       postgresRepos.Position,
		TrackingJob:    postgresRepos.TrackingJob,
		TrackingTask:   postgresRepos.TrackingTask,
		TrackingResult: postgresRepos.TrackingResult,
	}
}
