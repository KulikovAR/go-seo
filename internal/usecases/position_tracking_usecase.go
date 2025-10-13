package usecases

import (
	"fmt"
	"sync"
	"time"

	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/services"
)

type PositionTrackingUseCase struct {
	siteRepo     repositories.SiteRepository
	keywordRepo  repositories.KeywordRepository
	positionRepo repositories.PositionRepository
	xmlRiver     *services.XMLRiverService
}

func NewPositionTrackingUseCase(
	siteRepo repositories.SiteRepository,
	keywordRepo repositories.KeywordRepository,
	positionRepo repositories.PositionRepository,
	xmlRiver *services.XMLRiverService,
) *PositionTrackingUseCase {
	return &PositionTrackingUseCase{
		siteRepo:     siteRepo,
		keywordRepo:  keywordRepo,
		positionRepo: positionRepo,
		xmlRiver:     xmlRiver,
	}
}

func (uc *PositionTrackingUseCase) TrackSitePositions(siteID int, source, device, os string, ads bool, country, lang string, pages int) (int, error) {
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

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int
	var firstError error

	for _, keyword := range keywords {
		wg.Add(1)
		go func(kw *entities.Keyword) {
			defer wg.Done()

			err := uc.trackKeywordPosition(site, kw, source, device, os, ads, country, lang, pages)

			mu.Lock()
			if err != nil && firstError == nil {
				firstError = err
			} else if err == nil {
				count++
			}
			mu.Unlock()
		}(keyword)
	}

	wg.Wait()

	if firstError != nil {
		return count, firstError
	}

	return count, nil
}

func (uc *PositionTrackingUseCase) trackKeywordPosition(
	site *entities.Site,
	keyword *entities.Keyword,
	source, device, os string,
	ads bool,
	country, lang string,
	pages int,
) error {
	position, url, title, err := uc.xmlRiver.FindSitePosition(keyword.Value, site.Domain, source, pages, device, os, ads, country, lang)
	if err != nil {
		return &DomainError{
			Code:    ErrorPositionCreation,
			Message: fmt.Sprintf("Failed to search position for keyword '%s'", keyword.Value),
			Err:     err,
		}
	}

	positionEntity := &entities.Position{
		KeywordID: keyword.ID,
		SiteID:    site.ID,
		Rank:      position,
		URL:       url,
		Title:     title,
		Source:    source,
		Device:    device,
		OS:        os,
		Ads:       ads,
		Country:   country,
		Lang:      lang,
		Pages:     pages,
		Date:      time.Now(),
	}

	if err := uc.positionRepo.Create(positionEntity); err != nil {
		return &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to save position",
			Err:     err,
		}
	}

	return nil
}

func (uc *PositionTrackingUseCase) GetPositionsHistory(siteID int, keywordID *int) ([]*entities.Position, error) {
	var positions []*entities.Position
	var err error

	if keywordID != nil {
		positions, err = uc.positionRepo.GetByKeywordAndSite(*keywordID, siteID)
	} else {
		positions, err = uc.positionRepo.GetBySiteID(siteID)
	}

	if err != nil {
		return nil, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Failed to fetch positions history",
			Err:     err,
		}
	}

	for _, pos := range positions {
		keyword, err := uc.keywordRepo.GetByID(pos.KeywordID)
		if err == nil {
			pos.Keyword = keyword
		}
	}

	return positions, nil
}

func (uc *PositionTrackingUseCase) GetLatestPositions() ([]*entities.Position, error) {
	sites, err := uc.siteRepo.GetAll()
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Failed to fetch sites",
			Err:     err,
		}
	}

	var latestPositions []*entities.Position

	for _, site := range sites {
		keywords, err := uc.keywordRepo.GetBySiteID(site.ID)
		if err != nil {
			continue // Пропускаем сайты с ошибками
		}

		for _, keyword := range keywords {
			latestPosition, err := uc.positionRepo.GetLatestByKeywordAndSite(keyword.ID, site.ID)
			if err == nil && latestPosition != nil {
				latestPositions = append(latestPositions, latestPosition)
			}
		}
	}

	return latestPositions, nil
}
