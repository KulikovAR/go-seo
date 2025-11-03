package repositories

import "go-seo/internal/domain/entities"

type GroupRepository interface {
	Create(group *entities.Group) error
	GetByID(id int) (*entities.Group, error)
	GetAllBySite(siteID int) ([]*entities.Group, error)
	Update(group *entities.Group) error
	Delete(id int) error
}
