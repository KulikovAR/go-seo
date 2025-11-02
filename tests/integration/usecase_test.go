package integration

import (
	"testing"
	"time"

	"go-seo/internal/domain/entities"
	"go-seo/internal/usecases"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSiteRepository struct {
	mock.Mock
}

func (m *MockSiteRepository) Create(site *entities.Site) error {
	args := m.Called(site)
	site.ID = 1
	return args.Error(0)
}

func (m *MockSiteRepository) GetByID(id int) (*entities.Site, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Site), args.Error(1)
}

func (m *MockSiteRepository) GetByDomain(domain string) (*entities.Site, error) {
	args := m.Called(domain)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Site), args.Error(1)
}

func (m *MockSiteRepository) GetAll() ([]*entities.Site, error) {
	args := m.Called()
	return args.Get(0).([]*entities.Site), args.Error(1)
}

func (m *MockSiteRepository) GetByIDs(ids []int) ([]*entities.Site, error) {
	args := m.Called(ids)
	return args.Get(0).([]*entities.Site), args.Error(1)
}

func (m *MockSiteRepository) Update(site *entities.Site) error {
	args := m.Called(site)
	return args.Error(0)
}

func (m *MockSiteRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockKeywordRepository struct {
	mock.Mock
}

func (m *MockKeywordRepository) Create(keyword *entities.Keyword) error {
	args := m.Called(keyword)
	keyword.ID = 1 // Симулируем создание
	return args.Error(0)
}

func (m *MockKeywordRepository) GetByID(id int) (*entities.Keyword, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Keyword), args.Error(1)
}

func (m *MockKeywordRepository) GetByValueAndSite(value string, siteID int) (*entities.Keyword, error) {
	args := m.Called(value, siteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Keyword), args.Error(1)
}

func (m *MockKeywordRepository) GetBySiteID(siteID int) ([]*entities.Keyword, error) {
	args := m.Called(siteID)
	return args.Get(0).([]*entities.Keyword), args.Error(1)
}

func (m *MockKeywordRepository) GetAll() ([]*entities.Keyword, error) {
	args := m.Called()
	return args.Get(0).([]*entities.Keyword), args.Error(1)
}

func (m *MockKeywordRepository) Update(keyword *entities.Keyword) error {
	args := m.Called(keyword)
	return args.Error(0)
}

func (m *MockKeywordRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockKeywordRepository) GetByIDs(ids []int) ([]*entities.Keyword, error) {
	args := m.Called(ids)
	return args.Get(0).([]*entities.Keyword), args.Error(1)
}

func (m *MockKeywordRepository) CountBySiteID(siteID int) (int, error) {
	args := m.Called(siteID)
	return args.Get(0).(int), args.Error(1)
}

type MockPositionRepository struct {
	mock.Mock
}

func (m *MockPositionRepository) Create(position *entities.Position) error {
	args := m.Called(position)
	position.ID = 1
	return args.Error(0)
}

func (m *MockPositionRepository) CreateBatch(positions []*entities.Position) error {
	args := m.Called(positions)
	for i := range positions {
		positions[i].ID = i + 1
	}
	return args.Error(0)
}

func (m *MockPositionRepository) GetByID(id int) (*entities.Position, error) {
	args := m.Called(id)
	return args.Get(0).(*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetByKeywordAndSite(keywordID, siteID int) ([]*entities.Position, error) {
	args := m.Called(keywordID, siteID)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetBySiteID(siteID int) ([]*entities.Position, error) {
	args := m.Called(siteID)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetLatestByKeywordAndSite(keywordID, siteID int) (*entities.Position, error) {
	args := m.Called(keywordID, siteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetAll() ([]*entities.Position, error) {
	args := m.Called()
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) Update(position *entities.Position) error {
	args := m.Called(position)
	return args.Error(0)
}

func (m *MockPositionRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPositionRepository) DeleteBySiteID(siteID int) error {
	args := m.Called(siteID)
	return args.Error(0)
}

func (m *MockPositionRepository) DeleteByKeywordID(keywordID int) error {
	args := m.Called(keywordID)
	return args.Error(0)
}
func (m *MockPositionRepository) GetBySiteIDAndSource(siteID int, source string) ([]*entities.Position, error) {
	args := m.Called(siteID, source)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetByKeywordAndSiteAndSource(keywordID, siteID int, source string) ([]*entities.Position, error) {
	args := m.Called(keywordID, siteID, source)
	return args.Get(0).([]*entities.Position), args.Error(1)
}
func (m *MockPositionRepository) GetBySiteIDWithDateRange(siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	args := m.Called(siteID, dateFrom, dateTo)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetBySiteIDAndSourceWithDateRange(siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	args := m.Called(siteID, source, dateFrom, dateTo)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetByKeywordAndSiteWithDateRange(keywordID, siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	args := m.Called(keywordID, siteID, dateFrom, dateTo)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetByKeywordAndSiteAndSourceWithDateRange(keywordID, siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	args := m.Called(keywordID, siteID, source, dateFrom, dateTo)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetLastUpdateDateBySiteIDExcludingSource(siteID int, excludeSource string) (*time.Time, error) {
	args := m.Called(siteID, excludeSource)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*time.Time), args.Error(1)
}

func (m *MockPositionRepository) GetTodayByKeywordAndSiteAndSource(keywordID, siteID int, source string) (*entities.Position, error) {
	args := m.Called(keywordID, siteID, source)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) CreateOrUpdateToday(position *entities.Position) error {
	args := m.Called(position)
	return args.Error(0)
}

func (m *MockPositionRepository) GetHistoryBySiteIDWithOnePerDay(siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	args := m.Called(siteID, dateFrom, dateTo)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetHistoryBySiteIDAndSourceWithOnePerDay(siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	args := m.Called(siteID, source, dateFrom, dateTo)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetHistoryByKeywordAndSiteWithOnePerDay(keywordID, siteID int, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	args := m.Called(keywordID, siteID, dateFrom, dateTo)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetHistoryByKeywordAndSiteAndSourceWithOnePerDay(keywordID, siteID int, source string, dateFrom, dateTo *time.Time) ([]*entities.Position, error) {
	args := m.Called(keywordID, siteID, source, dateFrom, dateTo)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetLatestBySiteID(siteID int) ([]*entities.Position, error) {
	args := m.Called(siteID)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetLatestBySiteIDAndSource(siteID int, source string) ([]*entities.Position, error) {
	args := m.Called(siteID, source)
	return args.Get(0).([]*entities.Position), args.Error(1)
}

func (m *MockPositionRepository) GetPositionStatistics(siteID int, source string, dateFrom, dateTo time.Time) (*entities.PositionStatistics, error) {
	args := m.Called(siteID, source, dateFrom, dateTo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.PositionStatistics), args.Error(1)
}

func (m *MockPositionRepository) GetPositionsHistoryPaginated(siteID int, keywordID *int, source *string, dateFrom, dateTo *time.Time, last bool, page, perPage int) ([]*entities.Position, int64, error) {
	args := m.Called(siteID, keywordID, source, dateFrom, dateTo, last, page, perPage)
	return args.Get(0).([]*entities.Position), args.Get(1).(int64), args.Error(2)
}

func (m *MockPositionRepository) GetCombinedPositionsPaginated(siteID int, source *string, includeWordstat bool, wordstatSort bool, dateFrom, dateTo, dateSort *time.Time, sortType string, rankFrom, rankTo *int, page, perPage int) ([]*entities.CombinedPosition, int64, error) {
	args := m.Called(siteID, source, includeWordstat, wordstatSort, dateFrom, dateTo, dateSort, sortType, rankFrom, rankTo, page, perPage)
	return args.Get(0).([]*entities.CombinedPosition), args.Get(1).(int64), args.Error(2)
}

func TestSiteUseCase_CreateSite(t *testing.T) {
	mockSiteRepo := new(MockSiteRepository)
	mockPositionRepo := new(MockPositionRepository)
	mockKeywordRepo := new(MockKeywordRepository)

	useCase := usecases.NewSiteUseCase(mockSiteRepo, mockPositionRepo, mockKeywordRepo)

	mockSiteRepo.On("Create", mock.AnythingOfType("*entities.Site")).Return(nil)

	site, err := useCase.CreateSite("test.com")

	assert.NoError(t, err)
	assert.Equal(t, "test.com", site.Domain)
	assert.Equal(t, 1, site.ID)

	mockSiteRepo.AssertExpectations(t)
}

func TestKeywordUseCase_CreateKeyword(t *testing.T) {
	mockKeywordRepo := new(MockKeywordRepository)
	mockPositionRepo := new(MockPositionRepository)

	useCase := usecases.NewKeywordUseCase(mockKeywordRepo, mockPositionRepo)

	mockKeywordRepo.On("GetByValueAndSite", "купить чай", 1).Return(nil, assert.AnError)
	mockKeywordRepo.On("Create", mock.AnythingOfType("*entities.Keyword")).Return(nil)

	keyword, err := useCase.CreateKeyword("купить чай", 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, "купить чай", keyword.Value)
	assert.Equal(t, 1, keyword.SiteID)
	assert.Equal(t, 1, keyword.ID)

	mockKeywordRepo.AssertExpectations(t)
}

func TestKeywordUseCase_CreateKeyword_AlreadyExists(t *testing.T) {
	mockKeywordRepo := new(MockKeywordRepository)
	mockPositionRepo := new(MockPositionRepository)

	useCase := usecases.NewKeywordUseCase(mockKeywordRepo, mockPositionRepo)

	// Настраиваем мок - ключевое слово уже существует
	existingKeyword := &entities.Keyword{ID: 1, Value: "купить чай", SiteID: 1, GroupID: 1}
	mockKeywordRepo.On("GetByValueAndSite", "купить чай", 1).Return(existingKeyword, nil)

	keyword, err := useCase.CreateKeyword("купить чай", 1, 1)

	assert.Error(t, err)
	assert.Nil(t, keyword)
	assert.Contains(t, err.Error(), "Keyword already exists for this site")

	mockKeywordRepo.AssertExpectations(t)
}

func TestKeywordUseCase_GetKeywordsBySite(t *testing.T) {
	mockKeywordRepo := new(MockKeywordRepository)
	mockPositionRepo := new(MockPositionRepository)

	useCase := usecases.NewKeywordUseCase(mockKeywordRepo, mockPositionRepo)

	keywords := []*entities.Keyword{
		{ID: 1, Value: "купить чай", SiteID: 1, GroupID: 1},
		{ID: 2, Value: "купить кофе", SiteID: 1, GroupID: 1},
	}
	mockKeywordRepo.On("GetBySiteID", 1).Return(keywords, nil)

	result, err := useCase.GetKeywordsBySite(1)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "купить чай", result[0].Value)
	assert.Equal(t, "купить кофе", result[1].Value)

	mockKeywordRepo.AssertExpectations(t)
}
