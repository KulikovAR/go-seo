package usecases

import (
	"go-seo/internal/domain/entities"
	"time"
)

type SiteUseCaseInterface interface {
	CreateSite(domain string) (*entities.Site, error)
	DeleteSite(id int) error
	GetAllSites() ([]*entities.Site, error)
	GetSitesByIDs(ids []int) ([]*entities.Site, error)
	GetKeywordsCount(siteID int) (int, error)
	GetLastPositionUpdateDate(siteID int) (*time.Time, error)
}

type KeywordUseCaseInterface interface {
	CreateKeyword(value string, siteID int, groupID *int) (*entities.Keyword, error)
	CreateKeywordsBatch(keywords []*entities.Keyword) ([]*entities.Keyword, []error)
	UpdateKeyword(id int, groupID *int) (*entities.Keyword, error)
	DeleteKeyword(id int) error
	GetKeywordsBySite(siteID int) ([]*entities.Keyword, error)
}

type GroupUseCaseInterface interface {
	CreateGroup(name string, siteID int) (*entities.Group, error)
	UpdateGroup(id int, name string) (*entities.Group, error)
	DeleteGroup(id int) error
	GetGroupsBySite(siteID int) ([]*entities.Group, error)
}
