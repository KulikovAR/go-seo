package repositories

import (
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database"
	"go-seo/internal/infrastructure/database/postgres/models"

	"gorm.io/gorm"
)

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) repositories.GroupRepository {
	return &groupRepository{db: db}
}

func (r *groupRepository) Create(group *entities.Group) error {
	model := &models.Group{
		Name: group.Name,
	}

	if err := r.db.Create(model).Error; err != nil {
		return database.WrapDatabaseError(err)
	}

	group.ID = model.ID
	return nil
}

func (r *groupRepository) GetByID(id int) (*entities.Group, error) {
	var model models.Group
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *groupRepository) GetAll() ([]*entities.Group, error) {
	var models []models.Group
	if err := r.db.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	groups := make([]*entities.Group, len(models))
	for i, model := range models {
		groups[i] = r.toDomain(&model)
	}

	return groups, nil
}

func (r *groupRepository) Update(group *entities.Group) error {
	model := &models.Group{
		ID:   group.ID,
		Name: group.Name,
	}

	return r.db.Save(model).Error
}

func (r *groupRepository) Delete(id int) error {
	return r.db.Delete(&models.Group{}, id).Error
}

func (r *groupRepository) toDomain(model *models.Group) *entities.Group {
	return &entities.Group{
		ID:   model.ID,
		Name: model.Name,
	}
}
