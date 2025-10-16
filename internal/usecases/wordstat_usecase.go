package usecases

import (
	"fmt"
	"time"

	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/services"
)

type WordstatUseCase struct {
	siteRepo     repositories.SiteRepository
	keywordRepo  repositories.KeywordRepository
	positionRepo repositories.PositionRepository
	wordstat     *services.WordstatService
}

func NewWordstatUseCase(
	siteRepo repositories.SiteRepository,
	keywordRepo repositories.KeywordRepository,
	positionRepo repositories.PositionRepository,
	wordstat *services.WordstatService,
) *WordstatUseCase {
	return &WordstatUseCase{
		siteRepo:     siteRepo,
		keywordRepo:  keywordRepo,
		positionRepo: positionRepo,
		wordstat:     wordstat,
	}
}

func (uc *WordstatUseCase) TrackKeywordFrequency(keywordID int) (int, error) {
	keyword, err := uc.keywordRepo.GetByID(keywordID)
	if err != nil {
		return 0, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Keyword not found",
			Err:     err,
		}
	}

	frequency, err := uc.wordstat.GetKeywordFrequency(keyword.Value)
	if err != nil {
		return 0, &DomainError{
			Code:    ErrorPositionCreation,
			Message: fmt.Sprintf("Failed to get frequency for keyword '%s'", keyword.Value),
			Err:     err,
		}
	}

	positionEntity := &entities.Position{
		KeywordID: keyword.ID,
		SiteID:    keyword.SiteID,
		Rank:      frequency,
		URL:       "",
		Title:     "",
		Source:    entities.Wordstat,
		Device:    "",
		OS:        "",
		Ads:       false,
		Country:   "",
		Lang:      "",
		Pages:     0,
		Date:      time.Now(),
	}

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return 0, &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to save wordstat position",
			Err:     err,
		}
	}

	return frequency, nil
}

func (uc *WordstatUseCase) TrackSiteKeywordsFrequency(siteID int) (int, error) {
	site, err := uc.siteRepo.GetByID(siteID)
	if err != nil {
		return 0, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Site not found",
			Err:     err,
		}
	}

	keywords, err := uc.keywordRepo.GetBySiteID(siteID)
	if err != nil {
		return 0, &DomainError{
			Code:    ErrorPositionFetch,
			Message: fmt.Sprintf("Failed to fetch keywords for site %s", site.Domain),
			Err:     err,
		}
	}

	var count int
	for _, keyword := range keywords {
		_, err := uc.TrackKeywordFrequency(keyword.ID)
		if err == nil {
			count++
		}
	}

	return count, nil
}

func (uc *WordstatUseCase) GetRelatedKeywords(query string) ([]services.WordstatItem, error) {
	relatedKeywords, err := uc.wordstat.GetRelatedKeywords(query)
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorPositionFetch,
			Message: fmt.Sprintf("Failed to get related keywords for '%s'", query),
			Err:     err,
		}
	}

	return relatedKeywords, nil
}
