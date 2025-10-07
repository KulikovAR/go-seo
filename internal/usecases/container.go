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

func NewContainer(repos *repositories.Container, xmlRiver *services.XMLRiverService) *Container {
	return &Container{
		Site:             NewSiteUseCase(repos.Site, repos.Position),
		Keyword:          NewKeywordUseCase(repos.Keyword, repos.Position),
		PositionTracking: NewPositionTrackingUseCase(repos.Site, repos.Keyword, repos.Position, xmlRiver),
	}
}
