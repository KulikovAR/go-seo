package repositories

import (
	"go-seo/internal/domain/repositories"

	"gorm.io/gorm"
)

type RepositoryContainer struct {
	Keyword  repositories.KeywordRepository
	Site     repositories.SiteRepository
	Position repositories.PositionRepository
}

func NewRepositoryContainer(db *gorm.DB) *RepositoryContainer {
	return &RepositoryContainer{
		Keyword:  NewKeywordRepository(db),
		Site:     NewSiteRepository(db),
		Position: NewPositionRepository(db),
	}
}
