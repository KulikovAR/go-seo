package repositories

import (
	"go-seo/internal/domain/entities"
	"go-seo/internal/domain/repositories"
	positionModels "go-seo/internal/infrastructure/database/postgres/models"
	"sort"
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
	model := &positionModels.Position{
		KeywordID:         position.KeywordID,
		SiteID:            position.SiteID,
		Rank:              position.Rank,
		URL:               position.URL,
		Title:             position.Title,
		Source:            position.Source,
		Device:            position.Device,
		OS:                position.OS,
		Ads:               position.Ads,
		Country:           position.Country,
		Lang:              position.Lang,
		Pages:             position.Pages,
		Date:              position.Date,
		FilterGroupID:     position.FilterGroupID,
		WordstatQueryType: position.WordstatQueryType,
	}

	if err := r.db.Select("keyword_id", "site_id", "rank", "url", "title", "source", "device", "os", "ads", "country", "lang", "pages", "date", "filter_group_id", "wordstat_query_type").Create(model).Error; err != nil {
		return err
	}

	position.ID = model.ID
	return nil
}

func (r *positionRepository) CreateBatch(positions []*entities.Position) error {
	if len(positions) == 0 {
		return nil
	}

	models := make([]*positionModels.Position, len(positions))
	for i, position := range positions {
		models[i] = &positionModels.Position{
			KeywordID:         position.KeywordID,
			SiteID:            position.SiteID,
			Rank:              position.Rank,
			URL:               position.URL,
			Title:             position.Title,
			Source:            position.Source,
			Device:            position.Device,
			OS:                position.OS,
			Ads:               position.Ads,
			Country:           position.Country,
			Lang:              position.Lang,
			Pages:             position.Pages,
			Date:              position.Date,
			FilterGroupID:     position.FilterGroupID,
			WordstatQueryType: position.WordstatQueryType,
		}
	}

	if err := r.db.CreateInBatches(models, 100).Error; err != nil {
		return err
	}

	for i, model := range models {
		positions[i].ID = model.ID
	}

	return nil
}

func (r *positionRepository) GetByID(id int) (*entities.Position, error) {
	var model positionModels.Position
	if err := r.db.Preload("Keyword").Preload("Site").First(&model, id).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *positionRepository) GetByKeywordAndSite(keywordID, siteID int) ([]*entities.Position, error) {
	var models []positionModels.Position
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
	var models []positionModels.Position
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
	var models []positionModels.Position
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
	var models []positionModels.Position
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

	var models []positionModels.Position
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

	var models []positionModels.Position
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

	var models []positionModels.Position
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

	var models []positionModels.Position
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
	var model positionModels.Position
	if err := r.db.Where("keyword_id = ? AND site_id = ?", keywordID, siteID).
		Order("date DESC").
		First(&model).Error; err != nil {
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *positionRepository) GetAll() ([]*entities.Position, error) {
	var models []positionModels.Position
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
	return r.db.Model(&positionModels.Position{}).
		Where("id = ?", position.ID).
		Updates(positionModels.Position{
			KeywordID:         position.KeywordID,
			SiteID:            position.SiteID,
			Rank:              position.Rank,
			URL:               position.URL,
			Title:             position.Title,
			Source:            position.Source,
			Device:            position.Device,
			OS:                position.OS,
			Ads:               position.Ads,
			Country:           position.Country,
			Lang:              position.Lang,
			Pages:             position.Pages,
			Date:              position.Date,
			FilterGroupID:     position.FilterGroupID,
			WordstatQueryType: position.WordstatQueryType,
		}).Error
}

func (r *positionRepository) Delete(id int) error {
	return r.db.Delete(&positionModels.Position{}, id).Error
}

func (r *positionRepository) DeleteBySiteID(siteID int) error {
	return r.db.Where("site_id = ?", siteID).Delete(&positionModels.Position{}).Error
}

func (r *positionRepository) DeleteByKeywordID(keywordID int) error {
	return r.db.Where("keyword_id = ?", keywordID).Delete(&positionModels.Position{}).Error
}

func (r *positionRepository) GetTodayByKeywordAndSiteAndSource(keywordID, siteID int, source string, wordstatQueryType string, filterGroupID *int) (*entities.Position, error) {
	var model positionModels.Position

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	query := r.db.Where("keyword_id = ? AND site_id = ? AND source = ? AND date >= ? AND date <= ?",
		keywordID, siteID, source, startOfDay, endOfDay)

	if source == "wordstat" && wordstatQueryType != "" {
		query = query.Where("wordstat_query_type = ?", wordstatQueryType)
	}

	if (source == "google" || source == "yandex") && filterGroupID != nil {
		query = query.Where("filter_group_id = ?", *filterGroupID)
	} else if (source == "google" || source == "yandex") && filterGroupID == nil {
		query = query.Where("filter_group_id IS NULL")
	}

	if err := query.First(&model).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.toDomain(&model), nil
}

func (r *positionRepository) CreateOrUpdateToday(position *entities.Position) error {
	wordstatQueryType := ""
	if position.Source == "wordstat" {
		wordstatQueryType = position.WordstatQueryType
	}
	existingPosition, err := r.GetTodayByKeywordAndSiteAndSource(position.KeywordID, position.SiteID, position.Source, wordstatQueryType, position.FilterGroupID)
	if err != nil {
		return err
	}

	if existingPosition != nil {
		if position.Source == "google" || position.Source == "yandex" {
			existingFilterGroupID := existingPosition.FilterGroupID
			newFilterGroupID := position.FilterGroupID

			if (existingFilterGroupID == nil && newFilterGroupID == nil) ||
				(existingFilterGroupID != nil && newFilterGroupID != nil && *existingFilterGroupID == *newFilterGroupID) {
				position.ID = existingPosition.ID
				return r.Update(position)
			}
		} else {
			position.ID = existingPosition.ID
			return r.Update(position)
		}
	}

	return r.Create(position)
}

func (r *positionRepository) GetHistoryBySiteIDWithOnePerDay(siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	query := r.db.Table("positions").
		Select("DISTINCT ON (keyword_id, DATE(date)) *").
		Where("site_id = ?", siteID)

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	var models []positionModels.Position
	if err := query.Order("keyword_id, DATE(date) DESC, date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetHistoryBySiteIDAndSourceWithOnePerDay(siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	query := r.db.Table("positions").
		Select("DISTINCT ON (keyword_id, DATE(date)) *").
		Where("site_id = ? AND source = ?", siteID, source)

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	var models []positionModels.Position
	if err := query.Order("keyword_id, DATE(date) DESC, date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetHistoryByKeywordAndSiteWithOnePerDay(keywordID, siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	query := r.db.Table("positions").
		Select("DISTINCT ON (DATE(date)) *").
		Where("keyword_id = ? AND site_id = ?", keywordID, siteID)

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	var models []positionModels.Position
	if err := query.Order("DATE(date) DESC, date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetHistoryByKeywordAndSiteAndSourceWithOnePerDay(keywordID, siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	query := r.db.Table("positions").
		Select("DISTINCT ON (DATE(date)) *").
		Where("keyword_id = ? AND site_id = ? AND source = ?", keywordID, siteID, source)

	if dateFrom != nil {
		query = query.Where("date >= ?", *dateFrom)
	}
	if dateTo != nil {
		query = query.Where("date <= ?", *dateTo)
	}

	var models []positionModels.Position
	if err := query.Order("DATE(date) DESC, date DESC").Find(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetLatestBySiteID(siteID int) ([]*entities.Position, error) {
	var models []positionModels.Position

	query := `
		SELECT DISTINCT ON (keyword_id) *
		FROM positions 
		WHERE site_id = ?
		ORDER BY keyword_id, date DESC
	`

	if err := r.db.Raw(query, siteID).Scan(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) GetLatestBySiteIDAndSource(siteID int, source string) ([]*entities.Position, error) {
	var models []positionModels.Position

	query := `
		SELECT DISTINCT ON (keyword_id) *
		FROM positions 
		WHERE site_id = ? AND source = ?
		ORDER BY keyword_id, date DESC
	`

	if err := r.db.Raw(query, siteID, source).Scan(&models).Error; err != nil {
		return nil, err
	}

	positions := make([]*entities.Position, len(models))
	for i, model := range models {
		positions[i] = r.toDomain(&model)
	}

	return positions, nil
}

func (r *positionRepository) toDomain(model *positionModels.Position) *entities.Position {
	position := &entities.Position{
		ID:                model.ID,
		KeywordID:         model.KeywordID,
		SiteID:            model.SiteID,
		Rank:              model.Rank,
		URL:               model.URL,
		Title:             model.Title,
		Source:            model.Source,
		Device:            model.Device,
		OS:                model.OS,
		Ads:               model.Ads,
		Country:           model.Country,
		Lang:              model.Lang,
		Pages:             model.Pages,
		Date:              model.Date,
		FilterGroupID:     model.FilterGroupID,
		WordstatQueryType: model.WordstatQueryType,
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

func (r *positionRepository) GetPositionStatistics(siteID int, source string, dateFrom, dateTo time.Time, filterGroupID *int) (*entities.PositionStatistics, error) {
	var stats entities.PositionStatistics

	query := `
		SELECT 
			COUNT(*) as total_positions,
			COUNT(DISTINCT keyword_id) as keywords_count,
			COUNT(CASE WHEN rank > 0 THEN 1 END) as visible,
			COUNT(CASE WHEN rank = 0 THEN 1 END) as not_visible,
			ROUND(AVG(CASE WHEN rank > 0 THEN rank END), 2) as avg_position,
			MIN(CASE WHEN rank > 0 THEN rank END) as best_position,
			MAX(CASE WHEN rank > 0 THEN rank END) as worst_position,
			COUNT(CASE WHEN rank BETWEEN 1 AND 3 THEN 1 END) as top_3,
			COUNT(CASE WHEN rank BETWEEN 1 AND 10 THEN 1 END) as top_10,
			COUNT(CASE WHEN rank BETWEEN 1 AND 20 THEN 1 END) as top_20,
			COUNT(CASE WHEN rank = 0 THEN 1 END) as not_found,
			COUNT(CASE WHEN rank BETWEEN 1 AND 3 THEN 1 END) as range_1_3_count,
			COUNT(CASE WHEN rank BETWEEN 4 AND 10 THEN 1 END) as range_4_10_count,
			COUNT(CASE WHEN rank BETWEEN 11 AND 30 THEN 1 END) as range_11_30_count,
			COUNT(CASE WHEN rank BETWEEN 31 AND 50 THEN 1 END) as range_31_50_count,
			COUNT(CASE WHEN rank BETWEEN 51 AND 100 THEN 1 END) as range_51_100_count,
			COUNT(CASE WHEN rank > 100 THEN 1 END) as range_100_plus_count
		FROM positions 
		WHERE site_id = $1 
		  AND source = $2 
		  AND date >= $3::date 
		  AND date <= $4::date
	`

	queryParams := []interface{}{siteID, source, dateFrom, dateTo}
	if filterGroupID != nil {
		query += ` AND filter_group_id = $5`
		queryParams = append(queryParams, *filterGroupID)
	}

	var result struct {
		TotalPositions int     `json:"total_positions"`
		KeywordsCount  int     `json:"keywords_count"`
		Visible        int     `json:"visible"`
		NotVisible     int     `json:"not_visible"`
		AvgPosition    float64 `json:"avg_position"`
		BestPosition   int     `json:"best_position"`
		WorstPosition  int     `json:"worst_position"`
		Top3           int     `json:"top_3"`
		Top10          int     `json:"top_10"`
		Top20          int     `json:"top_20"`
		NotFound       int     `json:"not_found"`
		Range1_3       int     `json:"range_1_3_count"`
		Range4_10      int     `json:"range_4_10_count"`
		Range11_30     int     `json:"range_11_30_count"`
		Range31_50     int     `json:"range_31_50_count"`
		Range51_100    int     `json:"range_51_100_count"`
		Range100Plus   int     `json:"range_100_plus_count"`
	}

	rows, err := r.db.Raw(query, queryParams...).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Читаем результаты вручную
	if rows.Next() {
		err := rows.Scan(
			&result.TotalPositions,
			&result.KeywordsCount,
			&result.Visible,
			&result.NotVisible,
			&result.AvgPosition,
			&result.BestPosition,
			&result.WorstPosition,
			&result.Top3,
			&result.Top10,
			&result.Top20,
			&result.NotFound,
			&result.Range1_3,
			&result.Range4_10,
			&result.Range11_30,
			&result.Range31_50,
			&result.Range51_100,
			&result.Range100Plus,
		)
		if err != nil {
			return nil, err
		}
	}

	var medianPosition float64
	medianQuery := `
		SELECT PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY rank) as median_position
		FROM positions 
		WHERE site_id = $1 AND source = $2 AND date >= $3::date AND date <= $4::date AND rank > 0
	`
	medianParams := []interface{}{siteID, source, dateFrom, dateTo}
	if filterGroupID != nil {
		medianQuery += ` AND filter_group_id = $5`
		medianParams = append(medianParams, *filterGroupID)
	}
	if err := r.db.Raw(medianQuery, medianParams...).Scan(&medianPosition).Error; err != nil {
		medianPosition = 0
	}

	var trends struct {
		Improved int `json:"improved"`
		Declined int `json:"declined"`
		Stable   int `json:"stable"`
	}

	trendsQuery := `
		WITH recent_data AS (
			SELECT keyword_id, rank, date
			FROM positions 
			WHERE site_id = $1 AND source = $2 
			  AND date >= $3::date AND date <= $4::date
			  AND date >= CURRENT_DATE - INTERVAL '30 days'
	`
	trendsParams := []interface{}{siteID, source, dateFrom, dateTo}
	if filterGroupID != nil {
		trendsQuery += ` AND filter_group_id = $5`
		trendsParams = append(trendsParams, *filterGroupID)
	}
	trendsQuery += `
		),
		keyword_changes AS (
			SELECT 
				keyword_id,
				(SELECT rank FROM recent_data rd2 
				 WHERE rd2.keyword_id = rd1.keyword_id 
				 ORDER BY rd2.date ASC LIMIT 1) as first_rank,
				(SELECT rank FROM recent_data rd3 
				 WHERE rd3.keyword_id = rd1.keyword_id 
				 ORDER BY rd3.date DESC LIMIT 1) as last_rank
			FROM recent_data rd1
			GROUP BY keyword_id
		)
		SELECT 
			COUNT(CASE WHEN last_rank < first_rank THEN 1 END) as improved,
			COUNT(CASE WHEN last_rank > first_rank THEN 1 END) as declined,
			COUNT(CASE WHEN last_rank = first_rank THEN 1 END) as stable
		FROM keyword_changes
		WHERE first_rank IS NOT NULL AND last_rank IS NOT NULL
	`

	if err := r.db.Raw(trendsQuery, trendsParams...).Scan(&trends).Error; err != nil {
		trends.Improved = 0
		trends.Declined = 0
		trends.Stable = 0
	}

	// Заполняем структуру статистики
	stats.TotalPositions = result.TotalPositions
	stats.KeywordsCount = result.KeywordsCount
	stats.Visible = result.Visible
	stats.NotVisible = result.NotVisible
	stats.VisibilityStats.AvgPosition = result.AvgPosition
	stats.VisibilityStats.BestPosition = result.BestPosition
	stats.VisibilityStats.WorstPosition = result.WorstPosition
	stats.VisibilityStats.MedianPosition = int(medianPosition)

	stats.PositionDistribution = entities.PositionDistribution{
		Top3:     result.Top3,
		Top10:    result.Top10,
		Top20:    result.Top20,
		NotFound: result.NotFound,
	}

	stats.PositionRanges = entities.PositionRanges{
		Range1_3:     result.Range1_3,
		Range4_10:    result.Range4_10,
		Range11_30:   result.Range11_30,
		Range31_50:   result.Range31_50,
		Range51_100:  result.Range51_100,
		Range100Plus: result.Range100Plus,
		NotFound:     result.NotFound,
	}

	stats.Trends = entities.Trends{
		Improved: trends.Improved,
		Declined: trends.Declined,
		Stable:   trends.Stable,
	}

	return &stats, nil
}

func (r *positionRepository) GetPositionsHistoryPaginated(siteID int, keywordID *int, source *string, dateFrom, dateTo *time.Time, last bool, page, perPage int) ([]*entities.Position, int64, error) {
	var positions []*entities.Position
	var total int64
	var err error

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 50
	}
	if perPage > 100 {
		perPage = 100
	}

	offset := (page - 1) * perPage

	if last {
		if keywordID != nil && source != nil {
			position, err := r.GetLatestByKeywordAndSite(*keywordID, siteID)
			if err != nil {
				return nil, 0, err
			}
			if position != nil && position.Source == *source {
				positions = []*entities.Position{position}
				total = 1
			} else {
				positions = []*entities.Position{}
				total = 0
			}
		} else if keywordID != nil {
			position, err := r.GetLatestByKeywordAndSite(*keywordID, siteID)
			if err != nil {
				return nil, 0, err
			}
			if position != nil {
				positions = []*entities.Position{position}
				total = 1
			} else {
				positions = []*entities.Position{}
				total = 0
			}
		} else if source != nil {
			positions, err = r.GetLatestBySiteIDAndSource(siteID, *source)
			if err != nil {
				return nil, 0, err
			}
			total = int64(len(positions))
		} else {
			positions, err = r.GetLatestBySiteID(siteID)
			if err != nil {
				return nil, 0, err
			}
			total = int64(len(positions))
		}
		var query *gorm.DB
		var countQuery *gorm.DB

		if keywordID != nil && source != nil {
			query = r.db.Where("keyword_id = ? AND site_id = ? AND source = ?", *keywordID, siteID, *source)
			countQuery = r.db.Model(&positionModels.Position{}).Where("keyword_id = ? AND site_id = ? AND source = ?", *keywordID, siteID, *source)
		} else if keywordID != nil {
			query = r.db.Where("keyword_id = ? AND site_id = ?", *keywordID, siteID)
			countQuery = r.db.Model(&positionModels.Position{}).Where("keyword_id = ? AND site_id = ?", *keywordID, siteID)
		} else if source != nil {
			query = r.db.Where("site_id = ? AND source = ?", siteID, *source)
			countQuery = r.db.Model(&positionModels.Position{}).Where("site_id = ? AND source = ?", siteID, *source)
		} else {
			query = r.db.Where("site_id = ?", siteID)
			countQuery = r.db.Model(&positionModels.Position{}).Where("site_id = ?", siteID)
		}

		if dateFrom != nil {
			query = query.Where("date >= ?", *dateFrom)
			countQuery = countQuery.Where("date >= ?", *dateFrom)
		}
		if dateTo != nil {
			query = query.Where("date <= ?", *dateTo)
			countQuery = countQuery.Where("date <= ?", *dateTo)
		}

		if err := countQuery.Count(&total).Error; err != nil {
			return nil, 0, err
		}

		var models []positionModels.Position
		if err := query.Preload("Keyword").
			Order("date DESC").
			Offset(offset).
			Limit(perPage).
			Find(&models).Error; err != nil {
			return nil, 0, err
		}

		positions = make([]*entities.Position, len(models))
		for i, model := range models {
			positions[i] = r.toDomain(&model)
		}
	}

	return positions, total, nil
}

func (r *positionRepository) GetCombinedPositionsPaginated(siteID int, source *string, includeWordstat bool, wordstatSort bool, dateFrom, dateTo, dateSort *time.Time, sortType string, rankFrom, rankTo *int, groupID *int, filterGroupID *int, wordstatQueryType *string, page, perPage int) ([]*entities.CombinedPosition, int64, error) {
	offset := (page - 1) * perPage

	var allKeywords []positionModels.Keyword
	query := r.db.Where("site_id = ?", siteID)

	if groupID != nil {
		query = query.Where("group_id = ?", *groupID)
	}

	var totalKeywords int64
	if err := query.Model(&positionModels.Keyword{}).Count(&totalKeywords).Error; err != nil {
		return nil, 0, err
	}

	var keywords []positionModels.Keyword

	if wordstatSort {
		if err := query.Order("id").Find(&allKeywords).Error; err != nil {
			return nil, 0, err
		}

		if len(allKeywords) == 0 {
			return []*entities.CombinedPosition{}, totalKeywords, nil
		}

		type keywordWithPosition struct {
			keyword  positionModels.Keyword
			position int
		}

		var keywordsWithPositions []keywordWithPosition

		for _, keyword := range allKeywords {
			keywordID := keyword.ID

			var position int = 0

			wordstatQuery := r.db.Where("site_id = ? AND keyword_id = ? AND source = ?", siteID, keywordID, "wordstat")

			if dateFrom != nil {
				wordstatQuery = wordstatQuery.Where("date >= ?", *dateFrom)
			}
			if dateTo != nil {
				wordstatQuery = wordstatQuery.Where("date <= ?", *dateTo)
			}

			if wordstatQueryType != nil {
				wordstatQuery = wordstatQuery.Where("wordstat_query_type = ?", *wordstatQueryType)
			} else {
				wordstatQuery = wordstatQuery.Where("wordstat_query_type = ?", "default")
			}

			var wordstatModel positionModels.Position
			if err := wordstatQuery.Order("date DESC").First(&wordstatModel).Error; err == nil {
				position = wordstatModel.Rank
			}

			keywordsWithPositions = append(keywordsWithPositions, keywordWithPosition{
				keyword:  keyword,
				position: position,
			})
		}

		sort.Slice(keywordsWithPositions, func(i, j int) bool {
			posI := keywordsWithPositions[i].position
			posJ := keywordsWithPositions[j].position

			if posI == 0 && posJ == 0 {
				return false
			}
			if posI == 0 {
				return false
			}
			if posJ == 0 {
				return true
			}

			if sortType == "asc" {
				return posI < posJ
			} else {
				return posI > posJ
			}
		})

		start := offset
		end := offset + perPage
		if start > len(keywordsWithPositions) {
			start = len(keywordsWithPositions)
		}
		if end > len(keywordsWithPositions) {
			end = len(keywordsWithPositions)
		}

		if start >= end {
			return []*entities.CombinedPosition{}, totalKeywords, nil
		}

		for i := start; i < end; i++ {
			keywords = append(keywords, keywordsWithPositions[i].keyword)
		}
	} else if dateSort != nil {
		// Получаем все keywords
		if err := query.Order("id").Find(&allKeywords).Error; err != nil {
			return nil, 0, err
		}

		if len(allKeywords) == 0 {
			return []*entities.CombinedPosition{}, totalKeywords, nil
		}

		type keywordWithPosition struct {
			keyword  positionModels.Keyword
			position int
		}

		var keywordsWithPositions []keywordWithPosition

		for _, keyword := range allKeywords {
			keywordID := keyword.ID

			var position int = 0

			positionQuery := r.db.Where("site_id = ? AND keyword_id = ? AND source != ? AND DATE(date) = ?",
				siteID, keywordID, "wordstat", dateSort.Format("2006-01-02"))

			if filterGroupID != nil {
				positionQuery = positionQuery.Where("filter_group_id = ?", *filterGroupID)
			}

			if source != nil {
				if *source == "google" {
					positionQuery = positionQuery.Where("source = ?", "google")
				} else if *source == "yandex" {
					positionQuery = positionQuery.Where("source = ?", "yandex")
				}

				var positionModel positionModels.Position
				if err := positionQuery.Order("date DESC").First(&positionModel).Error; err == nil {
					position = positionModel.Rank
				}
			} else {
				positionQuery = positionQuery.Where("source IN ?", []string{"google", "yandex"})
				var positions []positionModels.Position
				if err := positionQuery.Order("rank ASC").Find(&positions).Error; err == nil {
					if len(positions) > 0 {
						position = positions[0].Rank
					}
				}
			}

			keywordsWithPositions = append(keywordsWithPositions, keywordWithPosition{
				keyword:  keyword,
				position: position,
			})
		}

		sort.Slice(keywordsWithPositions, func(i, j int) bool {
			posI := keywordsWithPositions[i].position
			posJ := keywordsWithPositions[j].position

			if posI == 0 && posJ == 0 {
				return false
			}
			if posI == 0 {
				return false
			}
			if posJ == 0 {
				return true
			}

			if sortType == "asc" {
				return posI < posJ
			} else {
				return posI > posJ
			}
		})

		// Применяем пагинацию к отсортированным keywords
		start := offset
		end := offset + perPage
		if start > len(keywordsWithPositions) {
			start = len(keywordsWithPositions)
		}
		if end > len(keywordsWithPositions) {
			end = len(keywordsWithPositions)
		}

		if start >= end {
			return []*entities.CombinedPosition{}, totalKeywords, nil
		}

		for i := start; i < end; i++ {
			keywords = append(keywords, keywordsWithPositions[i].keyword)
		}
	} else {
		// Старая логика: просто пагинация по id
		if err := query.Order("id").Offset(offset).Limit(perPage).Find(&keywords).Error; err != nil {
			return nil, 0, err
		}

		if len(keywords) == 0 {
			return []*entities.CombinedPosition{}, totalKeywords, nil
		}
	}

	keywordMap := make(map[int]*entities.Keyword)
	for _, kw := range keywords {
		keywordMap[kw.ID] = &entities.Keyword{
			ID:     kw.ID,
			Value:  kw.Value,
			SiteID: kw.SiteID,
		}
	}

	var allCombinedPositions []*entities.CombinedPosition

	for _, keyword := range keywords {
		keywordID := keyword.ID

		var positions []positionModels.Position

		query := r.db.Where("site_id = ? AND keyword_id = ? AND source != ?", siteID, keywordID, "wordstat")

		if source != nil {
			if *source == "google" {
				query = query.Where("source = ?", "google")
			} else if *source == "yandex" {
				query = query.Where("source = ?", "yandex")
			}
		}

		if dateFrom != nil {
			query = query.Where("date >= ?", *dateFrom)
		}
		if dateTo != nil {
			query = query.Where("date <= ?", *dateTo)
		}

		if rankFrom != nil {
			query = query.Where("rank >= ?", *rankFrom)
		}
		if rankTo != nil {
			query = query.Where("rank <= ?", *rankTo)
		}

		if filterGroupID != nil {
			query = query.Where("filter_group_id = ?", *filterGroupID)
		}

		if err := query.Order("date DESC").Find(&positions).Error; err != nil {
			return nil, 0, err
		}

		var googleYandexPositions []*entities.Position
		for _, model := range positions {
			position := r.toDomain(&model)
			googleYandexPositions = append(googleYandexPositions, position)
		}

		var wordstatPosition *entities.Position

		if includeWordstat {
			var wordstatModel positionModels.Position
			wordstatQuery := r.db.Where("site_id = ? AND keyword_id = ? AND source = ?", siteID, keywordID, "wordstat")

			if dateFrom != nil {
				wordstatQuery = wordstatQuery.Where("date >= ?", *dateFrom)
			}
			if dateTo != nil {
				wordstatQuery = wordstatQuery.Where("date <= ?", *dateTo)
			}

			if wordstatQueryType != nil {
				wordstatQuery = wordstatQuery.Where("wordstat_query_type = ?", *wordstatQueryType)
			} else {
				wordstatQuery = wordstatQuery.Where("wordstat_query_type = ?", "default")
			}

			if err := wordstatQuery.Order("date DESC").First(&wordstatModel).Error; err == nil {
				wordstatPosition = r.toDomain(&wordstatModel)
			}
		}

		if len(googleYandexPositions) == 0 {
			continue
		}

		var latestDate time.Time
		if len(positions) > 0 {
			latestDate = positions[0].Date
		}

		combined := &entities.CombinedPosition{
			ID:        keywordID,
			SiteID:    siteID,
			KeywordID: keywordID,
			Keyword:   keywordMap[keywordID],
			Date:      latestDate,
			Positions: googleYandexPositions,
			Wordstat:  wordstatPosition,
		}

		allCombinedPositions = append(allCombinedPositions, combined)
	}

	return allCombinedPositions, totalKeywords, nil
}

func (r *positionRepository) GetLastUpdateDateBySiteIDExcludingSource(siteID int, excludeSource string) (*time.Time, error) {
	var model positionModels.Position

	err := r.db.Where("site_id = ? AND source != ?", siteID, excludeSource).
		Order("date DESC").
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &model.Date, nil
}
