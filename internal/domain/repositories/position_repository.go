package repositories

import (
	"go-seo/internal/domain/entities"
	"time"
)

type PositionRepository interface {
	Create(position *entities.Position) error
	GetByID(id int) (*entities.Position, error)
	GetByKeywordAndSite(keywordID, siteID int) ([]*entities.Position, error)
	GetBySiteID(siteID int) ([]*entities.Position, error)
	GetBySiteIDAndSource(siteID int, source string) ([]*entities.Position, error)
	GetByKeywordAndSiteAndSource(keywordID, siteID int, source string) ([]*entities.Position, error)
	GetBySiteIDWithDateRange(siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error)
	GetBySiteIDAndSourceWithDateRange(siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error)
	GetByKeywordAndSiteWithDateRange(keywordID, siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error)
	GetByKeywordAndSiteAndSourceWithDateRange(keywordID, siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error)
	GetLatestByKeywordAndSite(keywordID, siteID int) (*entities.Position, error)
	GetAll() ([]*entities.Position, error)
	Update(position *entities.Position) error
	Delete(id int) error

	DeleteBySiteID(siteID int) error
	DeleteByKeywordID(keywordID int) error

	GetTodayByKeywordAndSiteAndSource(keywordID, siteID int, source string) (*entities.Position, error)
	CreateOrUpdateToday(position *entities.Position) error

	GetHistoryBySiteIDWithOnePerDay(siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error)
	GetHistoryBySiteIDAndSourceWithOnePerDay(siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error)
	GetHistoryByKeywordAndSiteWithOnePerDay(keywordID, siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error)
	GetHistoryByKeywordAndSiteAndSourceWithOnePerDay(keywordID, siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error)

	// Методы для получения последних данных по каждому ключевому слову
	GetLatestBySiteID(siteID int) ([]*entities.Position, error)
	GetLatestBySiteIDAndSource(siteID int, source string) ([]*entities.Position, error)
}
