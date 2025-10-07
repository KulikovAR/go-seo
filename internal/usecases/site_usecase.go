package usecases

import (
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/database"
)

type SiteUseCase struct {
	siteRepo     repositories.SiteRepository
	positionRepo repositories.PositionRepository
}

func NewSiteUseCase(siteRepo repositories.SiteRepository, positionRepo repositories.PositionRepository) *SiteUseCase {
	return &SiteUseCase{
		siteRepo:     siteRepo,
		positionRepo: positionRepo,
	}
}

func (uc *SiteUseCase) CreateSite(name, domain string) (*entities.Site, error) {
	existingSite, err := uc.siteRepo.GetByDomain(domain)
	if err == nil && existingSite != nil {
		return nil, &DomainError{
			Code:    ErrorSiteExists,
			Message: "Site with this domain already exists",
		}
	}

	site := &entities.Site{
		Name:   name,
		Domain: domain,
	}

	if err := uc.siteRepo.Create(site); err != nil {
		// Проверяем тип ошибки
		if database.IsDatabaseError(err) {
			switch database.GetDatabaseErrorCode(err) {
			case "DUPLICATE_ENTRY":
				return nil, &DomainError{
					Code:    ErrorSiteExists,
					Message: "Site with this domain already exists",
					Err:     err,
				}
			default:
				return nil, &DomainError{
					Code:    ErrorSiteCreation,
					Message: "Failed to create site",
					Err:     err,
				}
			}
		}
		return nil, &DomainError{
			Code:    ErrorSiteCreation,
			Message: "Failed to create site",
			Err:     err,
		}
	}

	return site, nil
}

func (uc *SiteUseCase) DeleteSite(id int) error {
	_, err := uc.siteRepo.GetByID(id)
	if err != nil {
		return &DomainError{
			Code:    ErrorSiteNotFound,
			Message: "Site not found",
			Err:     err,
		}
	}

	if err := uc.positionRepo.DeleteBySiteID(id); err != nil {
		return &DomainError{
			Code:    ErrorPositionDeletion,
			Message: "Failed to delete site positions",
			Err:     err,
		}
	}

	if err := uc.siteRepo.Delete(id); err != nil {
		return &DomainError{
			Code:    ErrorSiteDeletion,
			Message: "Failed to delete site",
			Err:     err,
		}
	}

	return nil
}

func (uc *SiteUseCase) GetAllSites() ([]*entities.Site, error) {
	sites, err := uc.siteRepo.GetAll()
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorSiteFetch,
			Message: "Failed to fetch sites",
			Err:     err,
		}
	}

	return sites, nil
}
