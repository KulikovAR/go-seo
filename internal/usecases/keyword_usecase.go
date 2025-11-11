package usecases

import (
	"fmt"
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database"
)

type KeywordUseCase struct {
	keywordRepo  repositories.KeywordRepository
	positionRepo repositories.PositionRepository
}

func NewKeywordUseCase(keywordRepo repositories.KeywordRepository, positionRepo repositories.PositionRepository) *KeywordUseCase {
	return &KeywordUseCase{
		keywordRepo:  keywordRepo,
		positionRepo: positionRepo,
	}
}

func (uc *KeywordUseCase) CreateKeyword(value string, siteID int, groupID *int) (*entities.Keyword, error) {
	existingKeyword, err := uc.keywordRepo.GetByValueAndSite(value, siteID)
	if err == nil && existingKeyword != nil {
		return nil, &DomainError{
			Code:    ErrorKeywordExists,
			Message: "Keyword already exists for this site",
		}
	}

	keyword := &entities.Keyword{
		Value:   value,
		SiteID:  siteID,
		GroupID: groupID,
	}

	if err := uc.keywordRepo.Create(keyword); err != nil {
		// Проверяем тип ошибки
		if database.IsDatabaseError(err) {
			switch database.GetDatabaseErrorCode(err) {
			case "FOREIGN_KEY_VIOLATION":
				return nil, &DomainError{
					Code:    ErrorKeywordCreation,
					Message: "Site not found",
					Err:     err,
				}
			case "DUPLICATE_ENTRY":
				return nil, &DomainError{
					Code:    ErrorKeywordExists,
					Message: "Keyword already exists for this site",
					Err:     err,
				}
			default:
				return nil, &DomainError{
					Code:    ErrorKeywordCreation,
					Message: "Failed to create keyword",
					Err:     err,
				}
			}
		}
		return nil, &DomainError{
			Code:    ErrorKeywordCreation,
			Message: "Failed to create keyword",
			Err:     err,
		}
	}

	return keyword, nil
}

func (uc *KeywordUseCase) CreateKeywordsBatch(keywords []*entities.Keyword) ([]*entities.Keyword, []error) {
	if len(keywords) == 0 {
		return []*entities.Keyword{}, []error{}
	}

	var toCreate []*entities.Keyword
	var errors []error

	for i, keyword := range keywords {
		existingKeyword, err := uc.keywordRepo.GetByValueAndSite(keyword.Value, keyword.SiteID)
		if err == nil && existingKeyword != nil {
			errors = append(errors, &DomainError{
				Code:    ErrorKeywordExists,
				Message: fmt.Sprintf("Keyword '%s' already exists for site %d", keyword.Value, keyword.SiteID),
			})
			continue
		}
		toCreate = append(toCreate, keywords[i])
	}

	if len(toCreate) == 0 {
		return []*entities.Keyword{}, errors
	}

	if err := uc.keywordRepo.CreateBatch(toCreate); err != nil {
		if database.IsDatabaseError(err) {
			switch database.GetDatabaseErrorCode(err) {
			case "FOREIGN_KEY_VIOLATION":
				for _, keyword := range toCreate {
					errors = append(errors, &DomainError{
						Code:    ErrorKeywordCreation,
						Message: fmt.Sprintf("Site %d not found for keyword '%s'", keyword.SiteID, keyword.Value),
						Err:     err,
					})
				}
			case "DUPLICATE_ENTRY":
				for _, keyword := range toCreate {
					errors = append(errors, &DomainError{
						Code:    ErrorKeywordExists,
						Message: fmt.Sprintf("Keyword '%s' already exists for site %d", keyword.Value, keyword.SiteID),
						Err:     err,
					})
				}
			default:
				for _, keyword := range toCreate {
					errors = append(errors, &DomainError{
						Code:    ErrorKeywordCreation,
						Message: fmt.Sprintf("Failed to create keyword '%s'", keyword.Value),
						Err:     err,
					})
				}
			}
		} else {
			for _, keyword := range toCreate {
				errors = append(errors, &DomainError{
					Code:    ErrorKeywordCreation,
					Message: fmt.Sprintf("Failed to create keyword '%s'", keyword.Value),
					Err:     err,
				})
			}
		}
		return []*entities.Keyword{}, errors
	}

	return toCreate, errors
}

func (uc *KeywordUseCase) UpdateKeyword(id int, groupID *int) (*entities.Keyword, error) {
	keyword, err := uc.keywordRepo.GetByID(id)
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorKeywordNotFound,
			Message: "Keyword not found",
			Err:     err,
		}
	}

	keyword.GroupID = groupID
	if err := uc.keywordRepo.Update(keyword); err != nil {
		return nil, &DomainError{
			Code:    ErrorKeywordUpdate,
			Message: "Failed to update keyword",
			Err:     err,
		}
	}

	return keyword, nil
}

func (uc *KeywordUseCase) DeleteKeyword(id int) error {
	_, err := uc.keywordRepo.GetByID(id)
	if err != nil {
		return &DomainError{
			Code:    ErrorKeywordNotFound,
			Message: "Keyword not found",
			Err:     err,
		}
	}

	if err := uc.positionRepo.DeleteByKeywordID(id); err != nil {
		return &DomainError{
			Code:    ErrorPositionDeletion,
			Message: "Failed to delete keyword positions",
			Err:     err,
		}
	}

	if err := uc.keywordRepo.Delete(id); err != nil {
		return &DomainError{
			Code:    ErrorKeywordDeletion,
			Message: "Failed to delete keyword",
			Err:     err,
		}
	}

	return nil
}

func (uc *KeywordUseCase) GetKeywordsBySite(siteID int) ([]*entities.Keyword, error) {
	keywords, err := uc.keywordRepo.GetBySiteID(siteID)
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorKeywordFetch,
			Message: "Failed to fetch keywords",
			Err:     err,
		}
	}

	return keywords, nil
}
