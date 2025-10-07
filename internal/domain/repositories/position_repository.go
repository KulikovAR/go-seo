package repositories

import "go-seo/internal/domain/entities"

type PositionRepository interface {
	Create(position *entities.Position) error
	GetByID(id int) (*entities.Position, error)
	GetByKeywordAndSite(keywordID, siteID int) ([]*entities.Position, error)
	GetLatestByKeywordAndSite(keywordID, siteID int) (*entities.Position, error)
	GetAll() ([]*entities.Position, error)
	Update(position *entities.Position) error
	Delete(id int) error

	DeleteBySiteID(siteID int) error
	DeleteByKeywordID(keywordID int) error
}
