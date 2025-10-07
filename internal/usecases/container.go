package usecases

import (
	"go-seo/internal/repositories"
)

type Container struct {
	Site    *SiteUseCase
	Keyword *KeywordUseCase
}

func NewContainer(repos *repositories.Container) *Container {
	return &Container{
		Site:    NewSiteUseCase(repos.Site, repos.Position),
		Keyword: NewKeywordUseCase(repos.Keyword, repos.Position),
	}
}
