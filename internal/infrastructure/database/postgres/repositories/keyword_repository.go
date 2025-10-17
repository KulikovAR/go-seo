package repositories

import (
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database"
	"go-seo/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type keywordRepository struct {
	db *gorm.DB
}

func NewKeywordRepository(db *gorm.DB) repositories.KeywordRepository {
	return &keywordRepository{db: db}
}

func (r *keywordRepository) Create(keyword *entities.Keyword) error {
	model := &models.Keyword{
		Value:  keyword.Value,
		SiteID: keyword.SiteID,
	}

	if err := r.db.Create(model).Error; err != nil {
		return database.WrapDatabaseError(err)
	}

	keyword.ID = model.ID
	return nil
}

func (r *keywordRepository) GetByID(id int) (*entities.Keyword, error) {
	var model models.Keyword
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *keywordRepository) GetByValueAndSite(value string, siteID int) (*entities.Keyword, error) {
	var model models.Keyword
	if err := r.db.Where("value = ? AND site_id = ?", value, siteID).First(&model).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *keywordRepository) GetBySiteID(siteID int) ([]*entities.Keyword, error) {
	var models []models.Keyword
	if err := r.db.Where("site_id = ?", siteID).Find(&models).Error; err != nil {
		return nil, err
	}

	keywords := make([]*entities.Keyword, len(models))
	for i, model := range models {
		keywords[i] = r.toDomain(&model)
	}

	return keywords, nil
}

func (r *keywordRepository) GetAll() ([]*entities.Keyword, error) {
	var modelKeywords []models.Keyword
	if err := r.db.Find(&modelKeywords).Error; err != nil {
		return nil, err
	}

	keywords := make([]*entities.Keyword, len(modelKeywords))
	for i, model := range modelKeywords {
		keywords[i] = r.toDomain(&model)
	}

	return keywords, nil
}

func (r *keywordRepository) Update(keyword *entities.Keyword) error {
	model := &models.Keyword{
		ID:     keyword.ID,
		Value:  keyword.Value,
		SiteID: keyword.SiteID,
	}

	return r.db.Save(model).Error
}

func (r *keywordRepository) Delete(id int) error {
	return r.db.Delete(&models.Keyword{}, id).Error
}

func (r *keywordRepository) CountBySiteID(siteID int) (int, error) {
	var count int64
	if err := r.db.Model(&models.Keyword{}).Where("site_id = ?", siteID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *keywordRepository) toDomain(model *models.Keyword) *entities.Keyword {
	return &entities.Keyword{
		ID:     model.ID,
		Value:  model.Value,
		SiteID: model.SiteID,
	}
}
