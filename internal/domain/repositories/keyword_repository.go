package repositories

import "go-seo/internal/domain/entities"

type KeywordRepository interface {
	Create(keyword *entities.Keyword) error
	GetByID(id int) (*entities.Keyword, error)
	GetByValueAndSite(value string, siteID int) (*entities.Keyword, error)
	GetBySiteID(siteID int) ([]*entities.Keyword, error)
	GetAll() ([]*entities.Keyword, error)
	Update(keyword *entities.Keyword) error
	Delete(id int) error
}
