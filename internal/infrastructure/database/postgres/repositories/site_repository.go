package repositories

import (
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database"
	"go-seo/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type siteRepository struct {
	db *gorm.DB
}

func NewSiteRepository(db *gorm.DB) repositories.SiteRepository {
	return &siteRepository{db: db}
}

func (r *siteRepository) Create(site *entities.Site) error {
	model := &models.Site{
		Domain: site.Domain,
	}

	if err := r.db.Create(model).Error; err != nil {
		return database.WrapDatabaseError(err)
	}

	site.ID = model.ID
	return nil
}

func (r *siteRepository) GetByID(id int) (*entities.Site, error) {
	var model models.Site
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *siteRepository) GetByDomain(domain string) (*entities.Site, error) {
	var model models.Site
	if err := r.db.Where("domain = ?", domain).First(&model).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *siteRepository) GetAll() ([]*entities.Site, error) {
	var models []models.Site
	if err := r.db.Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	sites := make([]*entities.Site, len(models))
	for i, model := range models {
		sites[i] = r.toDomain(&model)
	}

	return sites, nil
}

func (r *siteRepository) GetByIDs(ids []int) ([]*entities.Site, error) {
	if len(ids) == 0 {
		return []*entities.Site{}, nil
	}

	var models []models.Site
	if err := r.db.Where("id IN ?", ids).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	sites := make([]*entities.Site, len(models))
	for i, model := range models {
		sites[i] = r.toDomain(&model)
	}

	return sites, nil
}

func (r *siteRepository) Update(site *entities.Site) error {
	model := &models.Site{
		ID:     site.ID,
		Domain: site.Domain,
	}

	return r.db.Save(model).Error
}

func (r *siteRepository) Delete(id int) error {
	return r.db.Delete(&models.Site{}, id).Error
}

func (r *siteRepository) toDomain(model *models.Site) *entities.Site {
	return &entities.Site{
		ID:     model.ID,
		Domain: model.Domain,
	}
}
