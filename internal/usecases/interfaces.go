package usecases

import "go-seo/internal/domain/entities"

// Интерфейсы для тестирования
type SiteUseCaseInterface interface {
	CreateSite(name, domain string) (*entities.Site, error)
	DeleteSite(id int) error
	GetAllSites() ([]*entities.Site, error)
}

type KeywordUseCaseInterface interface {
	CreateKeyword(value string, siteID int) (*entities.Keyword, error)
	DeleteKeyword(id int) error
	GetKeywordsBySite(siteID int) ([]*entities.Keyword, error)
}
