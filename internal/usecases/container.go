package usecases

import (
	"go-seo/internal/infrastructure/services"
	"go-seo/internal/repositories"
)

type Container struct {
	Site             *SiteUseCase
	Keyword          *KeywordUseCase
	PositionTracking *PositionTrackingUseCase
}

func NewContainer(repos *repositories.Container, xmlRiver *services.XMLRiverService, wordstat *services.WordstatService) *Container {
	return &Container{
		Site:             NewSiteUseCase(repos.Site, repos.Position, repos.Keyword),
		Keyword:          NewKeywordUseCase(repos.Keyword, repos.Position),
		PositionTracking: NewPositionTrackingUseCase(repos.Site, repos.Keyword, repos.Position, xmlRiver, wordstat),
	}
}
