package usecases

import (
	"context"
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
	xmlStock     *services.XMLRiverService
	wordstat     *services.WordstatService
}

func NewPositionTrackingUseCase(
	siteRepo repositories.SiteRepository,
	keywordRepo repositories.KeywordRepository,
	positionRepo repositories.PositionRepository,
	xmlRiver *services.XMLRiverService,
	xmlStock *services.XMLRiverService,
	wordstat *services.WordstatService,
) *PositionTrackingUseCase {
	return &PositionTrackingUseCase{
		siteRepo:     siteRepo,
		keywordRepo:  keywordRepo,
		positionRepo: positionRepo,
		xmlRiver:     xmlRiver,
		xmlStock:     xmlStock,
		wordstat:     wordstat,
	}
}

func (uc *PositionTrackingUseCase) TrackGooglePositions(
	siteID int, device, os string, ads bool, country, lang string, pages int, subdomains bool,
	xmlUserID, xmlAPIKey, xmlBaseURL, tbs string, filter, highlights, nfpr, loc, ai int, raw string,
) (int, error) {
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

			err := uc.trackGoogleKeywordPosition(site, kw, device, os, ads, country, lang, pages, subdomains,
				xmlUserID, xmlAPIKey, xmlBaseURL, tbs, filter, highlights, nfpr, loc, ai, raw)

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

func (uc *PositionTrackingUseCase) TrackYandexPositions(
	siteID int, device, os string, ads bool, country, lang string, pages int, subdomains bool,
	xmlUserID, xmlAPIKey, xmlBaseURL string, groupBy, filter, highlights, within, lr int, raw string, inIndex, strict int,
) (int, error) {
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

			err := uc.trackYandexKeywordPosition(site, kw, device, os, ads, country, lang, pages, subdomains,
				xmlUserID, xmlAPIKey, xmlBaseURL, groupBy, filter, highlights, within, lr, raw, inIndex, strict)

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

func (uc *PositionTrackingUseCase) TrackWordstatPositions(siteID int, xmlUserID, xmlAPIKey, xmlBaseURL string, regions *int) (int, error) {
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

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int
	var firstError error

	semaphore := make(chan struct{}, 10)

	for _, keyword := range keywords {
		wg.Add(1)
		go func(kw *entities.Keyword) {
			defer wg.Done()

			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				mu.Lock()
				if firstError == nil {
					firstError = ctx.Err()
				}
				mu.Unlock()
				return
			}

			err := uc.trackWordstatKeywordPosition(kw, xmlUserID, xmlAPIKey, xmlBaseURL, regions)

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
	subdomains bool,
) error {
	if source == entities.Wordstat {
		return uc.trackWordstatPosition(keyword)
	}

	position, url, title, err := uc.xmlRiver.FindSitePositionWithSubdomains(keyword.Value, site.Domain, source, pages, device, os, ads, country, lang, subdomains, 0)
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

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to save position",
			Err:     err,
		}
	}

	return nil
}

func (uc *PositionTrackingUseCase) GetPositionsHistory(siteID int, keywordID *int, source *string, dateFrom, dateTo *time.Time, last bool) ([]*entities.Position, error) {
	var positions []*entities.Position
	var err error

	if last {
		if keywordID != nil && source != nil {
			position, err := uc.positionRepo.GetLatestByKeywordAndSite(*keywordID, siteID)
			if err != nil {
				return nil, &DomainError{
					Code:    ErrorPositionFetch,
					Message: "Failed to fetch latest position",
					Err:     err,
				}
			}
			if position != nil && position.Source == *source {
				positions = []*entities.Position{position}
			} else {
				positions = []*entities.Position{}
			}
		} else if keywordID != nil {
			position, err := uc.positionRepo.GetLatestByKeywordAndSite(*keywordID, siteID)
			if err != nil {
				return nil, &DomainError{
					Code:    ErrorPositionFetch,
					Message: "Failed to fetch latest position",
					Err:     err,
				}
			}
			if position != nil {
				positions = []*entities.Position{position}
			} else {
				positions = []*entities.Position{}
			}
		} else if source != nil {
			positions, err = uc.positionRepo.GetLatestBySiteIDAndSource(siteID, *source)
		} else {
			positions, err = uc.positionRepo.GetLatestBySiteID(siteID)
		}
		if keywordID != nil && source != nil {
			positions, err = uc.positionRepo.GetHistoryByKeywordAndSiteAndSourceWithOnePerDay(*keywordID, siteID, *source, dateFrom, dateTo)
		} else if keywordID != nil {
			positions, err = uc.positionRepo.GetHistoryByKeywordAndSiteWithOnePerDay(*keywordID, siteID, dateFrom, dateTo)
		} else if source != nil {
			positions, err = uc.positionRepo.GetHistoryBySiteIDAndSourceWithOnePerDay(siteID, *source, dateFrom, dateTo)
		} else {
			positions, err = uc.positionRepo.GetHistoryBySiteIDWithOnePerDay(siteID, dateFrom, dateTo)
		}
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

func (uc *PositionTrackingUseCase) GetPositionsHistoryPaginated(siteID int, keywordID *int, source *string, dateFrom, dateTo *time.Time, last bool, page, perPage int) ([]*entities.Position, int64, error) {
	positions, total, err := uc.positionRepo.GetPositionsHistoryPaginated(siteID, keywordID, source, dateFrom, dateTo, last, page, perPage)
	if err != nil {
		return nil, 0, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Failed to fetch positions history",
			Err:     err,
		}
	}

	for _, pos := range positions {
		if pos.Keyword == nil {
			keyword, err := uc.keywordRepo.GetByID(pos.KeywordID)
			if err == nil {
				pos.Keyword = keyword
			}
		}
	}

	return positions, total, nil
}

func (uc *PositionTrackingUseCase) GetPositionStatistics(siteID int, source string, dateFrom, dateTo time.Time) (*entities.PositionStatistics, error) {
	site, err := uc.siteRepo.GetByID(siteID)
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Site not found",
			Err:     err,
		}
	}
	if site == nil {
		return nil, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Site not found",
			Err:     fmt.Errorf("site with ID %d not found", siteID),
		}
	}

	if source != "google" && source != "yandex" && source != "wordstat" {
		return nil, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Invalid source. Must be 'google', 'yandex' or 'wordstat'",
			Err:     fmt.Errorf("invalid source: %s", source),
		}
	}

	stats, err := uc.positionRepo.GetPositionStatistics(siteID, source, dateFrom, dateTo)
	if err != nil {
		return nil, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Failed to fetch position statistics",
			Err:     err,
		}
	}

	return stats, nil
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

func (uc *PositionTrackingUseCase) GetCombinedPositionsPaginated(siteID int, source *string, includeWordstat bool, dateFrom, dateTo *time.Time, page, perPage int) ([]*entities.CombinedPosition, int64, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 50
	}
	if perPage > 100 {
		perPage = 100
	}

	if source != nil {
		if *source != "google" && *source != "yandex" {
			return nil, 0, &DomainError{
				Code:    ErrorPositionFetch,
				Message: "source must be either 'google' or 'yandex'",
				Err:     fmt.Errorf("invalid source: %s", *source),
			}
		}
	}

	combinedPositions, total, err := uc.positionRepo.GetCombinedPositionsPaginated(siteID, source, includeWordstat, dateFrom, dateTo, page, perPage)
	if err != nil {
		return nil, 0, &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Failed to fetch combined positions",
			Err:     err,
		}
	}

	return combinedPositions, total, nil
}

func (uc *PositionTrackingUseCase) trackWordstatPosition(keyword *entities.Keyword) error {
	frequency, err := uc.wordstat.GetKeywordFrequency(keyword.Value, nil)
	if err != nil {
		return &DomainError{
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
		return &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to save wordstat position",
			Err:     err,
		}
	}

	return nil
}

func (uc *PositionTrackingUseCase) trackGoogleKeywordPosition(
	site *entities.Site,
	keyword *entities.Keyword,
	device, os string,
	ads bool,
	country, lang string,
	pages int,
	subdomains bool,
	xmlUserID, xmlAPIKey, xmlBaseURL, tbs string,
	filter, highlights, nfpr, loc, ai int,
	raw string,
) error {
	// Создаем временный XMLRiver сервис с кастомными настройками
	var xmlRiverService *services.XMLRiverService
	if xmlUserID != "" && xmlAPIKey != "" && xmlBaseURL != "" {
		var err error
		xmlRiverService, err = services.NewXMLRiverService(xmlBaseURL, xmlUserID, xmlAPIKey)
		if err != nil {
			return &DomainError{
				Code:    ErrorPositionCreation,
				Message: "Failed to create XMLRiver service with custom settings",
				Err:     err,
			}
		}
	} else {
		// По умолчанию используем XMLStock для Google
		xmlRiverService = uc.xmlStock
	}

	position, url, title, err := xmlRiverService.FindSitePositionWithSubdomains(keyword.Value, site.Domain, entities.GoogleSearch, pages, device, os, ads, country, lang, subdomains, 0)
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
		Source:    entities.GoogleSearch,
		Device:    device,
		OS:        os,
		Ads:       ads,
		Country:   country,
		Lang:      lang,
		Pages:     pages,
		Date:      time.Now(),
	}

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to save position",
			Err:     err,
		}
	}

	return nil
}

func (uc *PositionTrackingUseCase) trackYandexKeywordPosition(
	site *entities.Site,
	keyword *entities.Keyword,
	device, os string,
	ads bool,
	country, lang string,
	pages int,
	subdomains bool,
	xmlUserID, xmlAPIKey, xmlBaseURL string,
	groupBy, filter, highlights, within, lr int,
	raw string,
	inIndex, strict int,
) error {
	// Создаем временный XMLRiver сервис с кастомными настройками
	var xmlRiverService *services.XMLRiverService
	if xmlUserID != "" && xmlAPIKey != "" && xmlBaseURL != "" {
		var err error
		xmlRiverService, err = services.NewXMLRiverService(xmlBaseURL, xmlUserID, xmlAPIKey)
		if err != nil {
			return &DomainError{
				Code:    ErrorPositionCreation,
				Message: "Failed to create XMLRiver service with custom settings",
				Err:     err,
			}
		}
	} else {
		// По умолчанию используем XMLStock для Yandex
		xmlRiverService = uc.xmlStock
	}

	position, url, title, err := xmlRiverService.FindSitePositionWithSubdomains(keyword.Value, site.Domain, entities.YandexSearch, pages, device, os, ads, country, lang, subdomains, lr)
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
		Source:    entities.YandexSearch,
		Device:    device,
		OS:        os,
		Ads:       ads,
		Country:   country,
		Lang:      lang,
		Pages:     pages,
		Date:      time.Now(),
	}

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to save position",
			Err:     err,
		}
	}

	return nil
}

func (uc *PositionTrackingUseCase) trackWordstatKeywordPosition(
	keyword *entities.Keyword,
	xmlUserID, xmlAPIKey, xmlBaseURL string,
	regions *int,
) error {
	// Создаем временный Wordstat сервис с кастомными настройками
	var wordstatService *services.WordstatService
	if xmlUserID != "" && xmlAPIKey != "" && xmlBaseURL != "" {
		var err error
		wordstatService, err = services.NewWordstatService(xmlBaseURL, xmlUserID, xmlAPIKey)
		if err != nil {
			return &DomainError{
				Code:    ErrorPositionCreation,
				Message: "Failed to create Wordstat service with custom settings",
				Err:     err,
			}
		}
	} else {
		// По умолчанию используем XMLRiver для Wordstat
		wordstatService = uc.wordstat
	}

	frequency, err := wordstatService.GetKeywordFrequency(keyword.Value, regions)
	if err != nil {
		return &DomainError{
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
		return &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to save wordstat position",
			Err:     err,
		}
	}

	return nil
}
