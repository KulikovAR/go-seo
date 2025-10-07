package repositories

import "go-seo/internal/domain/entities"

type SiteRepository interface {
	Create(site *entities.Site) error
	GetByID(id int) (*entities.Site, error)
	GetByDomain(domain string) (*entities.Site, error)
	GetAll() ([]*entities.Site, error)
	Update(site *entities.Site) error
	Delete(id int) error
}