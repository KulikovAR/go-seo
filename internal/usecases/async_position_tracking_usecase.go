package usecases

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	"go-seo/internal/infrastructure/services"
)

type AsyncPositionTrackingUseCase struct {
	siteRepo       repositories.SiteRepository
	keywordRepo    repositories.KeywordRepository
	positionRepo   repositories.PositionRepository
	jobRepo        repositories.TrackingJobRepository
	taskRepo       repositories.TrackingTaskRepository
	resultRepo     repositories.TrackingResultRepository
	xmlRiver       *services.XMLRiverService
	xmlStock       *services.XMLRiverService
	wordstat       *services.WordstatService
	kafkaService   *services.KafkaService
	idGenerator    *services.IDGeneratorService
	retryService   *services.RetryService
	workerPool     chan struct{}
	batchSize      int
	xmlRiverSoftID string
	xmlStockSoftID string
}

func NewAsyncPositionTrackingUseCase(
	siteRepo repositories.SiteRepository,
	keywordRepo repositories.KeywordRepository,
	positionRepo repositories.PositionRepository,
	jobRepo repositories.TrackingJobRepository,
	taskRepo repositories.TrackingTaskRepository,
	resultRepo repositories.TrackingResultRepository,
	xmlRiver *services.XMLRiverService,
	xmlStock *services.XMLRiverService,
	wordstat *services.WordstatService,
	kafkaService *services.KafkaService,
	idGenerator *services.IDGeneratorService,
	retryService *services.RetryService,
	workerCount int,
	batchSize int,
	xmlRiverSoftID string,
	xmlStockSoftID string,
) *AsyncPositionTrackingUseCase {
	return &AsyncPositionTrackingUseCase{
		siteRepo:       siteRepo,
		keywordRepo:    keywordRepo,
		positionRepo:   positionRepo,
		jobRepo:        jobRepo,
		taskRepo:       taskRepo,
		resultRepo:     resultRepo,
		xmlRiver:       xmlRiver,
		xmlStock:       xmlStock,
		wordstat:       wordstat,
		kafkaService:   kafkaService,
		idGenerator:    idGenerator,
		retryService:   retryService,
		workerPool:     make(chan struct{}, workerCount),
		batchSize:      batchSize,
		xmlRiverSoftID: xmlRiverSoftID,
		xmlStockSoftID: xmlStockSoftID,
	}
}

func (uc *AsyncPositionTrackingUseCase) StartAsyncGoogleTracking(
	siteID int, device, os string, ads bool, country, lang string, pages int, subdomains bool,
	xmlUserID, xmlAPIKey, xmlBaseURL, tbs string, filter, highlights, nfpr, loc, ai int, raw string,
	lr int, domain int, filterGroupID *int,
) (string, error) {
	site, err := uc.siteRepo.GetByID(siteID)
	if err != nil {
		return "", &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Site not found",
			Err:     err,
		}
	}

	keywords, err := uc.keywordRepo.GetBySiteID(siteID)
	if err != nil {
		return "", &DomainError{
			Code:    ErrorPositionFetch,
			Message: fmt.Sprintf("Failed to fetch keywords for site %s", site.Domain),
			Err:     err,
		}
	}

	jobID := uc.idGenerator.GenerateJobID()
	job := &entities.TrackingJob{
		ID:             jobID,
		SiteID:         siteID,
		Source:         entities.GoogleSearch,
		Status:         entities.TaskStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		TotalTasks:     len(keywords),
		CompletedTasks: 0,
		FailedTasks:    0,
	}

	if err := uc.jobRepo.Create(job); err != nil {
		return "", &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to create tracking job",
			Err:     err,
		}
	}

	var tasks []*entities.TrackingTask
	for _, keyword := range keywords {
		taskID := uc.idGenerator.GenerateTaskID()
		task := &entities.TrackingTask{
			ID:            taskID,
			JobID:         jobID,
			KeywordID:     keyword.ID,
			SiteID:        siteID,
			Source:        entities.GoogleSearch,
			Status:        entities.TaskStatusPending,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			RetryCount:    0,
			MaxRetries:    5,
			Device:        device,
			OS:            os,
			Ads:           ads,
			Country:       country,
			Lang:          lang,
			Pages:         pages,
			Subdomains:    subdomains,
			XMLUserID:     xmlUserID,
			XMLAPIKey:     xmlAPIKey,
			XMLBaseURL:    xmlBaseURL,
			TBS:           tbs,
			Filter:        filter,
			Highlights:    highlights,
			NFPR:          nfpr,
			Loc:           loc,
			AI:            ai,
			Raw:           raw,
			LR:            lr,
			Domain:        domain,
			FilterGroupID: filterGroupID,
		}
		tasks = append(tasks, task)
	}

	for _, task := range tasks {
		if err := uc.taskRepo.Create(task); err != nil {
			return "", &DomainError{
				Code:    ErrorPositionCreation,
				Message: "Failed to create tracking task",
				Err:     err,
			}
		}
	}

	go uc.processJob(jobID)

	return jobID, nil
}

func (uc *AsyncPositionTrackingUseCase) StartAsyncYandexTracking(
	siteID int, device, os string, ads bool, country, lang string, pages int, subdomains bool,
	xmlUserID, xmlAPIKey, xmlBaseURL string, groupBy, filter, highlights, within, lr int, raw string, inIndex, strict int,
	organic bool, filterGroupID *int,
) (string, error) {
	site, err := uc.siteRepo.GetByID(siteID)
	if err != nil {
		return "", &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Site not found",
			Err:     err,
		}
	}

	keywords, err := uc.keywordRepo.GetBySiteID(siteID)
	if err != nil {
		return "", &DomainError{
			Code:    ErrorPositionFetch,
			Message: fmt.Sprintf("Failed to fetch keywords for site %s", site.Domain),
			Err:     err,
		}
	}

	jobID := uc.idGenerator.GenerateJobID()
	job := &entities.TrackingJob{
		ID:             jobID,
		SiteID:         siteID,
		Source:         entities.YandexSearch,
		Status:         entities.TaskStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		TotalTasks:     len(keywords),
		CompletedTasks: 0,
		FailedTasks:    0,
	}

	if err := uc.jobRepo.Create(job); err != nil {
		return "", &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to create tracking job",
			Err:     err,
		}
	}

	var tasks []*entities.TrackingTask
	for _, keyword := range keywords {
		taskID := uc.idGenerator.GenerateTaskID()
		task := &entities.TrackingTask{
			ID:            taskID,
			JobID:         jobID,
			KeywordID:     keyword.ID,
			SiteID:        siteID,
			Source:        entities.YandexSearch,
			Status:        entities.TaskStatusPending,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			RetryCount:    0,
			MaxRetries:    5,
			Device:        device,
			OS:            os,
			Ads:           ads,
			Country:       country,
			Lang:          lang,
			Pages:         pages,
			Subdomains:    subdomains,
			XMLUserID:     xmlUserID,
			XMLAPIKey:     xmlAPIKey,
			XMLBaseURL:    xmlBaseURL,
			GroupBy:       groupBy,
			Filter:        filter,
			Highlights:    highlights,
			Within:        within,
			LR:            lr,
			Raw:           raw,
			InIndex:       inIndex,
			Strict:        strict,
			Organic:       organic,
			FilterGroupID: filterGroupID,
		}
		tasks = append(tasks, task)
	}

	for _, task := range tasks {
		if err := uc.taskRepo.Create(task); err != nil {
			return "", &DomainError{
				Code:    ErrorPositionCreation,
				Message: "Failed to create tracking task",
				Err:     err,
			}
		}
	}

	go uc.processJob(jobID)

	return jobID, nil
}

func (uc *AsyncPositionTrackingUseCase) StartAsyncWordstatTracking(
	siteID int, xmlUserID, xmlAPIKey, xmlBaseURL string, regions *int,
	defaultQuery, quotes, quotesExclamationMarks, exclamationMarks bool,
) (string, error) {
	site, err := uc.siteRepo.GetByID(siteID)
	if err != nil {
		return "", &DomainError{
			Code:    ErrorPositionFetch,
			Message: "Site not found",
			Err:     err,
		}
	}

	keywords, err := uc.keywordRepo.GetBySiteID(siteID)
	if err != nil {
		return "", &DomainError{
			Code:    ErrorPositionFetch,
			Message: fmt.Sprintf("Failed to fetch keywords for site %s", site.Domain),
			Err:     err,
		}
	}

	queryTypes := []string{}
	if defaultQuery {
		queryTypes = append(queryTypes, "default")
	}
	if quotes {
		queryTypes = append(queryTypes, "quotes")
	}
	if quotesExclamationMarks {
		queryTypes = append(queryTypes, "quotes_exclamation_marks")
	}
	if exclamationMarks {
		queryTypes = append(queryTypes, "exclamation_marks")
	}

	if len(queryTypes) == 0 {
		return "", &DomainError{
			Code:    ErrorPositionCreation,
			Message: "At least one query type must be enabled",
			Err:     fmt.Errorf("no query types enabled"),
		}
	}

	totalTasks := len(keywords) * len(queryTypes)

	jobID := uc.idGenerator.GenerateJobID()
	job := &entities.TrackingJob{
		ID:             jobID,
		SiteID:         siteID,
		Source:         entities.Wordstat,
		Status:         entities.TaskStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		TotalTasks:     totalTasks,
		CompletedTasks: 0,
		FailedTasks:    0,
	}

	if err := uc.jobRepo.Create(job); err != nil {
		return "", &DomainError{
			Code:    ErrorPositionCreation,
			Message: "Failed to create tracking job",
			Err:     err,
		}
	}

	var tasks []*entities.TrackingTask
	for _, keyword := range keywords {
		for _, queryType := range queryTypes {
			taskID := uc.idGenerator.GenerateTaskID()
			task := &entities.TrackingTask{
				ID:                taskID,
				JobID:             jobID,
				KeywordID:         keyword.ID,
				SiteID:            siteID,
				Source:            entities.Wordstat,
				Status:            entities.TaskStatusPending,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
				RetryCount:        0,
				MaxRetries:        5,
				XMLUserID:         xmlUserID,
				XMLAPIKey:         xmlAPIKey,
				XMLBaseURL:        xmlBaseURL,
				Regions:           regions,
				WordstatQueryType: queryType,
			}
			tasks = append(tasks, task)
		}
	}

	for _, task := range tasks {
		if err := uc.taskRepo.Create(task); err != nil {
			return "", &DomainError{
				Code:    ErrorPositionCreation,
				Message: "Failed to create tracking task",
				Err:     err,
			}
		}
	}

	go uc.processJob(jobID)

	return jobID, nil
}

func (uc *AsyncPositionTrackingUseCase) processJob(jobID string) {
	job, err := uc.jobRepo.GetByID(jobID)
	if err != nil {
		uc.kafkaService.SendJobStatus(jobID, string(entities.TaskStatusFailed), err.Error())
		return
	}

	if job.Status == entities.TaskStatusCompleted || job.Status == entities.TaskStatusFailed {
		return
	}

	if job.Status == entities.TaskStatusRunning {
		return
	}

	job.Status = entities.TaskStatusRunning
	uc.jobRepo.UpdateStatus(jobID, entities.TaskStatusRunning)
	uc.kafkaService.SendJobStatus(jobID, string(entities.TaskStatusRunning), "", 0)

	progressDone := make(chan struct{})
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				j, err := uc.jobRepo.GetByID(jobID)
				if err != nil {
					continue
				}
				if j.Status == entities.TaskStatusCompleted || j.Status == entities.TaskStatusFailed {
					return
				}
				if j.TotalTasks > 0 {
					p := (j.CompletedTasks + j.FailedTasks) * 100 / j.TotalTasks
					uc.kafkaService.SendJobStatus(jobID, string(entities.TaskStatusRunning), "", p)
				}
			case <-progressDone:
				return
			}
		}
	}()

	tasks, err := uc.taskRepo.GetByJobID(jobID)
	if err != nil {
		job.Status = entities.TaskStatusFailed
		job.Error = err.Error()
		uc.jobRepo.Update(job)
		uc.kafkaService.SendJobStatus(jobID, string(entities.TaskStatusFailed), err.Error())
		return
	}

	batchSize := uc.calculateOptimalBatchSize(len(tasks))
	batches := uc.createBatches(tasks, batchSize)

	var wg sync.WaitGroup
	for _, batch := range batches {
		wg.Add(1)
		go func(batchTasks []*entities.TrackingTask) {
			defer wg.Done()
			uc.processBatch(batchTasks)
		}(batch)
	}

	wg.Wait()

	job, _ = uc.jobRepo.GetByID(jobID)
	if job.FailedTasks == job.TotalTasks {
		job.Status = entities.TaskStatusFailed
		job.Error = "All tasks failed"
	} else {
		job.Status = entities.TaskStatusCompleted
		job.CompletedAt = &[]time.Time{time.Now()}[0]
	}
	uc.jobRepo.Update(job)

	close(progressDone)
	if job.Status == entities.TaskStatusCompleted {
		uc.kafkaService.SendJobStatus(jobID, string(job.Status), job.Error, 100)
		if job.Source == entities.GoogleSearch || job.Source == entities.YandexSearch {
			uc.calculateAndUpdateDynamic(job.SiteID, job.Source)
		}
	} else {
		uc.kafkaService.SendJobStatus(jobID, string(job.Status), job.Error)
	}
}

func (uc *AsyncPositionTrackingUseCase) createBatches(tasks []*entities.TrackingTask, batchSize int) [][]*entities.TrackingTask {
	var batches [][]*entities.TrackingTask

	for i := 0; i < len(tasks); i += batchSize {
		end := i + batchSize
		if end > len(tasks) {
			end = len(tasks)
		}
		batches = append(batches, tasks[i:end])
	}

	return batches
}

func (uc *AsyncPositionTrackingUseCase) calculateOptimalBatchSize(totalTasks int) int {
	workerCount := cap(uc.workerPool)

	if totalTasks <= workerCount {
		return 1
	}

	optimalBatchSize := totalTasks / workerCount

	minBatchSize := 1
	maxBatchSize := uc.batchSize

	if optimalBatchSize < minBatchSize {
		return minBatchSize
	}
	if optimalBatchSize > maxBatchSize {
		return maxBatchSize
	}

	return optimalBatchSize
}

func (uc *AsyncPositionTrackingUseCase) processBatch(batchTasks []*entities.TrackingTask) {
	uc.workerPool <- struct{}{}
	defer func() { <-uc.workerPool }()

	if len(batchTasks) == 0 {
		return
	}

	siteID := batchTasks[0].SiteID
	site, err := uc.siteRepo.GetByID(siteID)
	if err != nil {
		for _, task := range batchTasks {
			task.Status = entities.TaskStatusFailed
			task.Error = fmt.Sprintf("Failed to load site: %v", err)
			uc.taskRepo.Update(task)
			uc.updateJobProgress(task.JobID, false)
		}
		return
	}

	keywordIDs := make([]int, len(batchTasks))
	for i, task := range batchTasks {
		keywordIDs[i] = task.KeywordID
	}

	keywords, err := uc.keywordRepo.GetByIDs(keywordIDs)
	if err != nil {
		for _, task := range batchTasks {
			task.Status = entities.TaskStatusFailed
			task.Error = fmt.Sprintf("Failed to load keywords: %v", err)
			uc.taskRepo.Update(task)
			uc.updateJobProgress(task.JobID, false)
		}
		return
	}

	keywordMap := make(map[int]*entities.Keyword)
	for _, keyword := range keywords {
		keywordMap[keyword.ID] = keyword
	}

	for _, task := range batchTasks {
		keyword, exists := keywordMap[task.KeywordID]
		if !exists {
			task.Status = entities.TaskStatusFailed
			task.Error = "Keyword not found"
			uc.taskRepo.Update(task)
			uc.updateJobProgress(task.JobID, false)
			continue
		}

		uc.processTaskWithData(task, site, keyword)
	}
}

func (uc *AsyncPositionTrackingUseCase) processTaskWithData(task *entities.TrackingTask, site *entities.Site, keyword *entities.Keyword) {
	if task.Status == entities.TaskStatusCompleted || task.Status == entities.TaskStatusFailed {
		return
	}

	task.Status = entities.TaskStatusRunning
	uc.taskRepo.UpdateStatus(task.ID, entities.TaskStatusRunning)

	err := uc.retryService.ExecuteWithRetry(func() error {
		return uc.executeTaskWithData(task, site, keyword)
	})

	if err != nil {
		task.Status = entities.TaskStatusFailed
		task.Error = err.Error()
		uc.taskRepo.Update(task)
		uc.kafkaService.SendTaskStatus(&services.TaskStatusMessage{
			TaskID:    task.ID,
			JobID:     task.JobID,
			Status:    string(entities.TaskStatusFailed),
			Timestamp: time.Now(),
			Error:     err.Error(),
		})
		uc.updateJobProgress(task.JobID, false)
		return
	}

	task.Status = entities.TaskStatusCompleted
	task.CompletedAt = &[]time.Time{time.Now()}[0]
	uc.taskRepo.Update(task)

	uc.kafkaService.SendTaskStatus(&services.TaskStatusMessage{
		TaskID:    task.ID,
		JobID:     task.JobID,
		Status:    string(entities.TaskStatusCompleted),
		Timestamp: time.Now(),
	})

	uc.updateJobProgress(task.JobID, true)
}

func (uc *AsyncPositionTrackingUseCase) processTask(task *entities.TrackingTask) {
	task.Status = entities.TaskStatusRunning
	uc.taskRepo.UpdateStatus(task.ID, entities.TaskStatusRunning)

	err := uc.retryService.ExecuteWithRetry(func() error {
		return uc.executeTask(task)
	})

	if err != nil {
		task.Status = entities.TaskStatusFailed
		task.Error = err.Error()
		uc.taskRepo.Update(task)
		uc.kafkaService.SendTaskStatus(&services.TaskStatusMessage{
			TaskID:    task.ID,
			JobID:     task.JobID,
			Status:    string(entities.TaskStatusFailed),
			Timestamp: time.Now(),
			Error:     err.Error(),
		})
		uc.updateJobProgress(task.JobID, false)
		return
	}

	task.Status = entities.TaskStatusCompleted
	task.CompletedAt = &[]time.Time{time.Now()}[0]
	uc.taskRepo.Update(task)

	uc.kafkaService.SendTaskStatus(&services.TaskStatusMessage{
		TaskID:    task.ID,
		JobID:     task.JobID,
		Status:    string(entities.TaskStatusCompleted),
		Timestamp: time.Now(),
	})

	uc.updateJobProgress(task.JobID, true)
}

func (uc *AsyncPositionTrackingUseCase) executeTaskWithData(task *entities.TrackingTask, site *entities.Site, keyword *entities.Keyword) error {
	switch task.Source {
	case entities.GoogleSearch:
		return uc.executeGoogleTaskWithData(task, site, keyword)
	case entities.YandexSearch:
		return uc.executeYandexTaskWithData(task, site, keyword)
	case entities.Wordstat:
		return uc.executeWordstatTaskWithData(task, site, keyword)
	default:
		return fmt.Errorf("unknown source: %s", task.Source)
	}
}

func (uc *AsyncPositionTrackingUseCase) executeTask(task *entities.TrackingTask) error {
	switch task.Source {
	case entities.GoogleSearch:
		return uc.executeGoogleTask(task)
	case entities.YandexSearch:
		return uc.executeYandexTask(task)
	case entities.Wordstat:
		return uc.executeWordstatTask(task)
	default:
		return fmt.Errorf("unknown source: %s", task.Source)
	}
}

func (uc *AsyncPositionTrackingUseCase) executeGoogleTaskWithData(task *entities.TrackingTask, site *entities.Site, keyword *entities.Keyword) error {

	var xmlRiverService *services.XMLRiverService
	if task.XMLUserID != "" && task.XMLAPIKey != "" && task.XMLBaseURL != "" {
		var err error
		softID := uc.getSoftIDByBaseURL(task.XMLBaseURL)
		xmlRiverService, err = services.NewXMLRiverService(task.XMLBaseURL, task.XMLUserID, task.XMLAPIKey, softID)
		if err != nil {
			return err
		}
	} else {
		xmlRiverService = uc.xmlStock
	}

	// Для Google используем organic=false и groupBy=0
	position, url, title, err := xmlRiverService.FindSitePositionWithSubdomains(
		keyword.Value, site.Domain, entities.GoogleSearch, task.Pages,
		task.Device, task.OS, task.Ads, task.Country, task.Lang, task.Subdomains, task.LR, task.Domain,
		false, 0,
	)
	if err != nil {
		return err
	}

	positionEntity := &entities.Position{
		KeywordID:     keyword.ID,
		SiteID:        site.ID,
		Rank:          position,
		URL:           url,
		Title:         title,
		Source:        entities.GoogleSearch,
		Device:        task.Device,
		OS:            task.OS,
		Ads:           task.Ads,
		Country:       task.Country,
		Lang:          task.Lang,
		Pages:         task.Pages,
		Date:          time.Now(),
		FilterGroupID: task.FilterGroupID,
	}

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return err
	}

	result := &entities.TrackingResult{
		TaskID:    task.ID,
		JobID:     task.JobID,
		KeywordID: keyword.ID,
		SiteID:    site.ID,
		Source:    entities.GoogleSearch,
		Rank:      position,
		URL:       url,
		Title:     title,
		Device:    task.Device,
		OS:        task.OS,
		Ads:       task.Ads,
		Country:   task.Country,
		Lang:      task.Lang,
		Pages:     task.Pages,
		Date:      time.Now(),
		Success:   true,
	}

	return uc.resultRepo.Create(result)
}

func (uc *AsyncPositionTrackingUseCase) getSoftIDByBaseURL(baseURL string) string {
	baseURLLower := strings.ToLower(baseURL)
	if strings.Contains(baseURLLower, "xmlriver") {
		return uc.xmlRiverSoftID
	}
	if strings.Contains(baseURLLower, "xmlstock") {
		return uc.xmlStockSoftID
	}
	return uc.xmlRiverSoftID
}

func (uc *AsyncPositionTrackingUseCase) executeGoogleTask(task *entities.TrackingTask) error {
	site, err := uc.siteRepo.GetByID(task.SiteID)
	if err != nil {
		return err
	}

	keyword, err := uc.keywordRepo.GetByID(task.KeywordID)
	if err != nil {
		return err
	}

	var xmlRiverService *services.XMLRiverService
	if task.XMLUserID != "" && task.XMLAPIKey != "" && task.XMLBaseURL != "" {
		var err error
		softID := uc.getSoftIDByBaseURL(task.XMLBaseURL)
		xmlRiverService, err = services.NewXMLRiverService(task.XMLBaseURL, task.XMLUserID, task.XMLAPIKey, softID)
		if err != nil {
			return err
		}
	} else {
		xmlRiverService = uc.xmlStock
	}

	// Для Google используем organic=false и groupBy=0
	position, url, title, err := xmlRiverService.FindSitePositionWithSubdomains(
		keyword.Value, site.Domain, entities.GoogleSearch, task.Pages,
		task.Device, task.OS, task.Ads, task.Country, task.Lang, task.Subdomains, task.LR, task.Domain,
		false, 0,
	)
	if err != nil {
		return err
	}

	positionEntity := &entities.Position{
		KeywordID: keyword.ID,
		SiteID:    site.ID,
		Rank:      position,
		URL:       url,
		Title:     title,
		Source:    entities.GoogleSearch,
		Device:    task.Device,
		OS:        task.OS,
		Ads:       task.Ads,
		Country:   task.Country,
		Lang:      task.Lang,
		Pages:     task.Pages,
		Date:      time.Now(),
	}

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return err
	}

	result := &entities.TrackingResult{
		TaskID:    task.ID,
		JobID:     task.JobID,
		KeywordID: keyword.ID,
		SiteID:    site.ID,
		Source:    entities.GoogleSearch,
		Rank:      position,
		URL:       url,
		Title:     title,
		Device:    task.Device,
		OS:        task.OS,
		Ads:       task.Ads,
		Country:   task.Country,
		Lang:      task.Lang,
		Pages:     task.Pages,
		Date:      time.Now(),
		Success:   true,
	}

	return uc.resultRepo.Create(result)
}

func (uc *AsyncPositionTrackingUseCase) executeYandexTaskWithData(task *entities.TrackingTask, site *entities.Site, keyword *entities.Keyword) error {
	var xmlRiverService *services.XMLRiverService
	if task.XMLUserID != "" && task.XMLAPIKey != "" && task.XMLBaseURL != "" {
		var err error
		softID := uc.getSoftIDByBaseURL(task.XMLBaseURL)
		xmlRiverService, err = services.NewXMLRiverService(task.XMLBaseURL, task.XMLUserID, task.XMLAPIKey, softID)
		if err != nil {
			return err
		}
	} else {
		xmlRiverService = uc.xmlStock
	}

	// Если organic=false, используем groupby=pages*10 для получения всех результатов сразу
	var groupBy int
	if !task.Organic && task.Pages > 0 {
		groupBy = task.Pages * 10
	} else {
		groupBy = task.GroupBy
	}

	position, url, title, err := xmlRiverService.FindSitePositionWithSubdomains(
		keyword.Value, site.Domain, entities.YandexSearch, task.Pages,
		task.Device, task.OS, task.Ads, task.Country, task.Lang, task.Subdomains, task.LR, 0,
		task.Organic, groupBy,
	)
	if err != nil {
		return err
	}

	positionEntity := &entities.Position{
		KeywordID:     keyword.ID,
		SiteID:        site.ID,
		Rank:          position,
		URL:           url,
		Title:         title,
		Source:        entities.YandexSearch,
		Device:        task.Device,
		OS:            task.OS,
		Ads:           task.Ads,
		Country:       task.Country,
		Lang:          task.Lang,
		Pages:         task.Pages,
		Date:          time.Now(),
		FilterGroupID: task.FilterGroupID,
	}

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return err
	}

	result := &entities.TrackingResult{
		TaskID:    task.ID,
		JobID:     task.JobID,
		KeywordID: keyword.ID,
		SiteID:    site.ID,
		Source:    entities.YandexSearch,
		Rank:      position,
		URL:       url,
		Title:     title,
		Device:    task.Device,
		OS:        task.OS,
		Ads:       task.Ads,
		Country:   task.Country,
		Lang:      task.Lang,
		Pages:     task.Pages,
		Date:      time.Now(),
		Success:   true,
	}

	return uc.resultRepo.Create(result)
}

func (uc *AsyncPositionTrackingUseCase) executeWordstatTaskWithData(task *entities.TrackingTask, site *entities.Site, keyword *entities.Keyword) error {
	var wordstatService *services.WordstatService
	if task.XMLUserID != "" && task.XMLAPIKey != "" && task.XMLBaseURL != "" {
		var err error
		wordstatService, err = services.NewWordstatService(task.XMLBaseURL, task.XMLUserID, task.XMLAPIKey)
		if err != nil {
			return err
		}
	} else {
		wordstatService = uc.wordstat
	}

	queryType := task.WordstatQueryType
	if queryType == "" {
		queryType = "default"
	}

	modifiedQuery := uc.modifyWordstatQuery(keyword.Value, queryType)
	frequency, err := wordstatService.GetKeywordFrequency(modifiedQuery, keyword.Value, task.Regions)
	if err != nil {
		return err
	}

	positionEntity := &entities.Position{
		KeywordID:         keyword.ID,
		SiteID:            keyword.SiteID,
		Rank:              frequency,
		URL:               "",
		Title:             "",
		Source:            entities.Wordstat,
		Device:            "",
		OS:                "",
		Ads:               false,
		Country:           "",
		Lang:              "",
		Pages:             0,
		Date:              time.Now(),
		WordstatQueryType: queryType,
	}

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return err
	}

	result := &entities.TrackingResult{
		TaskID:    task.ID,
		JobID:     task.JobID,
		KeywordID: keyword.ID,
		SiteID:    keyword.SiteID,
		Source:    entities.Wordstat,
		Rank:      frequency,
		URL:       "",
		Title:     "",
		Device:    "",
		OS:        "",
		Ads:       false,
		Country:   "",
		Lang:      "",
		Pages:     0,
		Date:      time.Now(),
		Success:   true,
	}

	return uc.resultRepo.Create(result)
}

func (uc *AsyncPositionTrackingUseCase) executeYandexTask(task *entities.TrackingTask) error {
	site, err := uc.siteRepo.GetByID(task.SiteID)
	if err != nil {
		return err
	}

	keyword, err := uc.keywordRepo.GetByID(task.KeywordID)
	if err != nil {
		return err
	}

	var xmlRiverService *services.XMLRiverService
	if task.XMLUserID != "" && task.XMLAPIKey != "" && task.XMLBaseURL != "" {
		var err error
		softID := uc.getSoftIDByBaseURL(task.XMLBaseURL)
		xmlRiverService, err = services.NewXMLRiverService(task.XMLBaseURL, task.XMLUserID, task.XMLAPIKey, softID)
		if err != nil {
			return err
		}
	} else {
		xmlRiverService = uc.xmlStock
	}

	// Если organic=false, используем groupby=pages*10 для получения всех результатов сразу
	var groupBy int
	if !task.Organic && task.Pages > 0 {
		groupBy = task.Pages * 10
	} else {
		groupBy = task.GroupBy
	}

	position, url, title, err := xmlRiverService.FindSitePositionWithSubdomains(
		keyword.Value, site.Domain, entities.YandexSearch, task.Pages,
		task.Device, task.OS, task.Ads, task.Country, task.Lang, task.Subdomains, task.LR, 0,
		task.Organic, groupBy,
	)
	if err != nil {
		return err
	}

	positionEntity := &entities.Position{
		KeywordID: keyword.ID,
		SiteID:    site.ID,
		Rank:      position,
		URL:       url,
		Title:     title,
		Source:    entities.YandexSearch,
		Device:    task.Device,
		OS:        task.OS,
		Ads:       task.Ads,
		Country:   task.Country,
		Lang:      task.Lang,
		Pages:     task.Pages,
		Date:      time.Now(),
	}

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return err
	}

	result := &entities.TrackingResult{
		TaskID:    task.ID,
		JobID:     task.JobID,
		KeywordID: keyword.ID,
		SiteID:    site.ID,
		Source:    entities.YandexSearch,
		Rank:      position,
		URL:       url,
		Title:     title,
		Device:    task.Device,
		OS:        task.OS,
		Ads:       task.Ads,
		Country:   task.Country,
		Lang:      task.Lang,
		Pages:     task.Pages,
		Date:      time.Now(),
		Success:   true,
	}

	return uc.resultRepo.Create(result)
}

func (uc *AsyncPositionTrackingUseCase) modifyWordstatQuery(query string, queryType string) string {
	switch queryType {
	case "default":
		return query
	case "quotes":
		return fmt.Sprintf(`"%s"`, query)
	case "quotes_exclamation_marks":
		words := strings.Fields(query)
		modifiedWords := make([]string, len(words))
		for i, word := range words {
			modifiedWords[i] = "!" + word
		}
		return fmt.Sprintf(`"%s"`, strings.Join(modifiedWords, " "))
	case "exclamation_marks":
		words := strings.Fields(query)
		modifiedWords := make([]string, len(words))
		for i, word := range words {
			modifiedWords[i] = "!" + word
		}
		return fmt.Sprintf(`"[%s]"`, strings.Join(modifiedWords, " "))
	default:
		return query
	}
}

func (uc *AsyncPositionTrackingUseCase) executeWordstatTask(task *entities.TrackingTask) error {
	keyword, err := uc.keywordRepo.GetByID(task.KeywordID)
	if err != nil {
		return err
	}

	var wordstatService *services.WordstatService
	if task.XMLUserID != "" && task.XMLAPIKey != "" && task.XMLBaseURL != "" {
		var err error
		wordstatService, err = services.NewWordstatService(task.XMLBaseURL, task.XMLUserID, task.XMLAPIKey)
		if err != nil {
			return err
		}
	} else {
		wordstatService = uc.wordstat
	}

	queryType := task.WordstatQueryType
	if queryType == "" {
		queryType = "default"
	}

	modifiedQuery := uc.modifyWordstatQuery(keyword.Value, queryType)
	frequency, err := wordstatService.GetKeywordFrequency(modifiedQuery, keyword.Value, task.Regions)
	if err != nil {
		return err
	}

	positionEntity := &entities.Position{
		KeywordID:         keyword.ID,
		SiteID:            keyword.SiteID,
		Rank:              frequency,
		URL:               "",
		Title:             "",
		Source:            entities.Wordstat,
		Device:            "",
		OS:                "",
		Ads:               false,
		Country:           "",
		Lang:              "",
		Pages:             0,
		Date:              time.Now(),
		WordstatQueryType: queryType,
	}

	if err := uc.positionRepo.CreateOrUpdateToday(positionEntity); err != nil {
		return err
	}

	result := &entities.TrackingResult{
		TaskID:    task.ID,
		JobID:     task.JobID,
		KeywordID: keyword.ID,
		SiteID:    keyword.SiteID,
		Source:    entities.Wordstat,
		Rank:      frequency,
		URL:       "",
		Title:     "",
		Device:    "",
		OS:        "",
		Ads:       false,
		Country:   "",
		Lang:      "",
		Pages:     0,
		Date:      time.Now(),
		Success:   true,
	}

	return uc.resultRepo.Create(result)
}

func (uc *AsyncPositionTrackingUseCase) updateJobProgress(jobID string, success bool) {
	job, err := uc.jobRepo.GetByID(jobID)
	if err != nil {
		return
	}

	if success {
		job.CompletedTasks++
	} else {
		job.FailedTasks++
	}

	uc.jobRepo.UpdateProgress(jobID, job.CompletedTasks, job.FailedTasks)

	if job.TotalTasks > 0 {
		progress := (job.CompletedTasks + job.FailedTasks) * 100 / job.TotalTasks
		uc.kafkaService.SendJobStatus(jobID, string(entities.TaskStatusRunning), "", progress)
	}
}

func (uc *AsyncPositionTrackingUseCase) calculateAndUpdateDynamic(siteID int, source string) {
	currentPositions, err := uc.positionRepo.GetLatestBySiteIDAndSource(siteID, source)
	if err != nil {
		return
	}

	type currentPositionData struct {
		rank int
		date time.Time
	}

	currentMap := make(map[int]currentPositionData)
	for _, pos := range currentPositions {
		if pos.Rank > 0 {
			currentMap[pos.KeywordID] = currentPositionData{
				rank: pos.Rank,
				date: pos.Date,
			}
		}
	}

	if len(currentMap) == 0 {
		site, err := uc.siteRepo.GetByID(siteID)
		if err != nil {
			return
		}
		if source == entities.GoogleSearch {
			site.GoogleDynamic = nil
		} else if source == entities.YandexSearch {
			site.YandexDynamic = nil
		}
		uc.siteRepo.Update(site)
		return
	}

	var totalDiff int
	hasComparisons := false
	for kwID, currentData := range currentMap {
		dateBefore := currentData.date.Add(-time.Nanosecond)
		positions, err := uc.positionRepo.GetByKeywordAndSiteAndSourceWithDateRange(kwID, siteID, source, nil, &dateBefore)
		if err != nil {
			continue
		}

		if len(positions) == 0 {
			continue
		}

		var previousRank int
		for _, pos := range positions {
			if pos.Rank > 0 {
				previousRank = pos.Rank
				break
			}
		}

		if previousRank > 0 {
			diff := previousRank - currentData.rank
			totalDiff += diff
			hasComparisons = true
		}
	}

	site, err := uc.siteRepo.GetByID(siteID)
	if err != nil {
		return
	}

	var dynamic *int
	if hasComparisons {
		if totalDiff > 0 {
			val := 1
			dynamic = &val
		} else if totalDiff < 0 {
			val := 0
			dynamic = &val
		}
	}

	if source == entities.GoogleSearch {
		site.GoogleDynamic = dynamic
	} else if source == entities.YandexSearch {
		site.YandexDynamic = dynamic
	}

	uc.siteRepo.Update(site)
}
