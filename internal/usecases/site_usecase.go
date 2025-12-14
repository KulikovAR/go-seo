package usecases

import (
	"time"

	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
)

type SiteUseCase struct {
	siteRepo     repositories.SiteRepository
	positionRepo repositories.PositionRepository
	keywordRepo  repositories.KeywordRepository
	groupRepo    repositories.GroupRepository
	jobRepo      repositories.TrackingJobRepository
	taskRepo     repositories.TrackingTaskRepository
	resultRepo   repositories.TrackingResultRepository
}

func NewSiteUseCase(
	siteRepo repositories.SiteRepository,
	positionRepo repositories.PositionRepository,
	keywordRepo repositories.KeywordRepository,
	groupRepo repositories.GroupRepository,
	jobRepo repositories.TrackingJobRepository,
	taskRepo repositories.TrackingTaskRepository,
	resultRepo repositories.TrackingResultRepository,
) *SiteUseCase {
	return &SiteUseCase{
		siteRepo:     siteRepo,
		positionRepo: positionRepo,
		keywordRepo:  keywordRepo,
		groupRepo:    groupRepo,
		jobRepo:      jobRepo,
		taskRepo:     taskRepo,
		resultRepo:   resultRepo,
	}
}

func (uc *SiteUseCase) CreateSite(domain string) (*entities.Site, error) {
	site := &entities.Site{
		Domain: domain,
	}

	if err := uc.siteRepo.Create(site); err != nil {
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

	if err := uc.resultRepo.DeleteBySiteID(id); err != nil {
		return &DomainError{
			Code:    ErrorPositionDeletion,
			Message: "Failed to delete site tracking results",
			Err:     err,
		}
	}

	if err := uc.taskRepo.DeleteBySiteID(id); err != nil {
		return &DomainError{
			Code:    ErrorPositionDeletion,
			Message: "Failed to delete site tracking tasks",
			Err:     err,
		}
	}

	if err := uc.jobRepo.DeleteBySiteID(id); err != nil {
		return &DomainError{
			Code:    ErrorPositionDeletion,
			Message: "Failed to delete site tracking jobs",
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

	if err := uc.keywordRepo.DeleteBySiteID(id); err != nil {
		return &DomainError{
			Code:    ErrorPositionDeletion,
			Message: "Failed to delete site keywords",
			Err:     err,
		}
	}

	if err := uc.groupRepo.DeleteBySiteID(id); err != nil {
		return &DomainError{
			Code:    ErrorPositionDeletion,
			Message: "Failed to delete site groups",
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

func (uc *SiteUseCase) GetSitesByIDs(ids []int) ([]*entities.Site, error) {
	sites, err := uc.siteRepo.GetByIDs(ids)
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorSiteFetch,
			Message: "Failed to fetch sites by IDs",
			Err:     err,
		}
	}

	return sites, nil
}

func (uc *SiteUseCase) GetKeywordsCount(siteID int) (int, error) {
	count, err := uc.keywordRepo.CountBySiteID(siteID)
	if err != nil {
		return 0, &DomainError{
			Code:    ErrorSiteFetch,
			Message: "Failed to get keywords count",
			Err:     err,
		}
	}

	return count, nil
}

func (uc *SiteUseCase) GetLastPositionUpdateDate(siteID int) (*time.Time, error) {
	date, err := uc.positionRepo.GetLastUpdateDateBySiteIDExcludingSource(siteID, "wordstat")
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorSiteFetch,
			Message: "Failed to get last position update date",
			Err:     err,
		}
	}

	return date, nil
}
