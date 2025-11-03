package usecases

import (
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database"
)

type GroupUseCase struct {
	groupRepo repositories.GroupRepository
}

func NewGroupUseCase(groupRepo repositories.GroupRepository) *GroupUseCase {
	return &GroupUseCase{
		groupRepo: groupRepo,
	}
}

func (uc *GroupUseCase) CreateGroup(name string, siteID int) (*entities.Group, error) {
	group := &entities.Group{
		Name:   name,
		SiteID: siteID,
	}

	if err := uc.groupRepo.Create(group); err != nil {
		if database.IsDatabaseError(err) {
			switch database.GetDatabaseErrorCode(err) {
			case "DUPLICATE_ENTRY":
				return nil, &DomainError{
					Code:    ErrorGroupExists,
					Message: "Group with this name already exists",
					Err:     err,
				}
			default:
				return nil, &DomainError{
					Code:    ErrorGroupCreation,
					Message: "Failed to create group",
					Err:     err,
				}
			}
		}
		return nil, &DomainError{
			Code:    ErrorGroupCreation,
			Message: "Failed to create group",
			Err:     err,
		}
	}

	return group, nil
}

func (uc *GroupUseCase) UpdateGroup(id int, name string) (*entities.Group, error) {
	group, err := uc.groupRepo.GetByID(id)
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorGroupNotFound,
			Message: "Group not found",
			Err:     err,
		}
	}

	group.Name = name
	if err := uc.groupRepo.Update(group); err != nil {
		return nil, &DomainError{
			Code:    ErrorGroupFetch,
			Message: "Failed to update group",
			Err:     err,
		}
	}

	return group, nil
}

func (uc *GroupUseCase) DeleteGroup(id int) error {
	_, err := uc.groupRepo.GetByID(id)
	if err != nil {
		return &DomainError{
			Code:    ErrorGroupNotFound,
			Message: "Group not found",
			Err:     err,
		}
	}

	if err := uc.groupRepo.Delete(id); err != nil {
		return &DomainError{
			Code:    ErrorGroupDeletion,
			Message: "Failed to delete group",
			Err:     err,
		}
	}

	return nil
}

func (uc *GroupUseCase) GetGroupsBySite(siteID int) ([]*entities.Group, error) {
	groups, err := uc.groupRepo.GetAllBySite(siteID)
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorGroupFetch,
			Message: "Failed to fetch groups",
			Err:     err,
		}
	}

	return groups, nil
}
