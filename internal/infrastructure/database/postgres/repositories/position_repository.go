package repositories

import (
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database/postgres/models"
	"time"

	"gorm.io/gorm"
)

type positionRepository struct {
	db *gorm.DB
}

func NewPositionRepository(db *gorm.DB) repositories.PositionRepository {
	return &positionRepository{db: db}
}

func (r *positionRepository) Create(position *entities.Position) error {
	model := &models.Position{
		KeywordID: position.KeywordID,
		SiteID:    position.SiteID,
		Rank:      position.Rank,
		URL:       position.URL,
		Title:     position.Title,
		Source:    position.Source,
		Device:    position.Device,
		OS:        position.OS,
		Ads:       position.Ads,
		Country:   position.Country,
		Lang:      position.Lang,
		Pages:     position.Pages,
		Date:      position.Date,
	}

	if err := r.db.Create(model).Error; err != nil {
		return err
	}

	position.ID = model.ID
	return nil
}

func (r *positionRepository) GetByID(id int) (*entities.Position, error) {
	var model models.Position
	if err := r.db.Preload("Keyword").Preload("Site").First(&model, id).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *positionRepository) GetByKeywordAndSite(keywordID, siteID int) ([]*entities.Position, error) {
	var models []models.Position
	if err := r.db.Where("keyword_id = ? AND site_id = ?", keywordID, siteID).
		Order("date DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetBySiteID(siteID int) ([]*entities.Position, error) {
	var models []models.Position
	if err := r.db.Where("site_id = ?", siteID).
		Order("date DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}
func (r *positionRepository) GetBySiteIDAndSource(siteID int, source string) ([]*entities.Position, error) {
	var models []models.Position
	if err := r.db.Where("site_id = ? AND source = ?", siteID, source).
		Order("date DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetByKeywordAndSiteAndSource(keywordID, siteID int, source string) ([]*entities.Position, error) {
	var models []models.Position
	if err := r.db.Where("keyword_id = ? AND site_id = ? AND source = ?", keywordID, siteID, source).
		Order("date DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}
func (r *positionRepository) GetBySiteIDWithDateRange(siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	query := r.db.Where("site_id = ?", siteID)

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	var models []models.Position
	if err := query.Order("date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetBySiteIDAndSourceWithDateRange(siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	query := r.db.Where("site_id = ? AND source = ?", siteID, source)

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	var models []models.Position
	if err := query.Order("date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetByKeywordAndSiteWithDateRange(keywordID, siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	query := r.db.Where("keyword_id = ? AND site_id = ?", keywordID, siteID)

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	var models []models.Position
	if err := query.Order("date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetByKeywordAndSiteAndSourceWithDateRange(keywordID, siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	query := r.db.Where("keyword_id = ? AND site_id = ? AND source = ?", keywordID, siteID, source)

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	var models []models.Position
	if err := query.Order("date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetLatestByKeywordAndSite(keywordID, siteID int) (*entities.Position, error) {
	var model models.Position
	if err := r.db.Where("keyword_id = ? AND site_id = ?", keywordID, siteID).
		Order("date DESC").
		First(&model).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *positionRepository) GetAll() ([]*entities.Position, error) {
	var models []models.Position
	if err := r.db.Preload("Keyword").Preload("Site").
		Order("date DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) Update(position *entities.Position) error {
	model := &models.Position{
		ID:        position.ID,
		KeywordID: position.KeywordID,
		SiteID:    position.SiteID,
		Rank:      position.Rank,
		URL:       position.URL,
		Title:     position.Title,
		Source:    position.Source,
		Device:    position.Device,
		OS:        position.OS,
		Ads:       position.Ads,
		Country:   position.Country,
		Lang:      position.Lang,
		Pages:     position.Pages,
		Date:      position.Date,
	}

	return r.db.Save(model).Error
}

func (r *positionRepository) Delete(id int) error {
	return r.db.Delete(&models.Position{}, id).Error
}

func (r *positionRepository) DeleteBySiteID(siteID int) error {
	return r.db.Where("site_id = ?", siteID).Delete(&models.Position{}).Error
}

func (r *positionRepository) DeleteByKeywordID(keywordID int) error {
	return r.db.Where("keyword_id = ?", keywordID).Delete(&models.Position{}).Error
}

func (r *positionRepository) toDomain(model *models.Position) *entities.Position {
	position := &entities.Position{
		ID:        model.ID,
		KeywordID: model.KeywordID,
		SiteID:    model.SiteID,
		Rank:      model.Rank,
		URL:       model.URL,
		Title:     model.Title,
		Source:    model.Source,
		Device:    model.Device,
		OS:        model.OS,
		Ads:       model.Ads,
		Country:   model.Country,
		Lang:      model.Lang,
		Pages:     model.Pages,
		Date:      model.Date,
	}

	if model.Keyword.ID != 0 {
		position.Keyword = &entities.Keyword{
			ID:    model.Keyword.ID,
			Value: model.Keyword.Value,
		}
	}

	if model.Site.ID != 0 {
		position.Site = &entities.Site{
			ID:     model.Site.ID,
			Domain: model.Site.Domain,
		}
	}

	return position
}
