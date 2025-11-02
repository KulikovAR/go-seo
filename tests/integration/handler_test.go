package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-seo/internal/delivery/http/dto"
	"go-seo/internal/delivery/http/handlers"
	"go-seo/internal/domain/entities"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSiteHandler_CreateSite(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSiteUseCase := new(MockSiteUseCase)
	handler := handlers.NewSiteHandler(mockSiteUseCase)

	expectedSite := &entities.Site{ID: 1, Domain: "test.com"}
	mockSiteUseCase.On("CreateSite", "test.com").Return(expectedSite, nil)

	reqBody := dto.CreateSiteRequest{
		Domain: "test.com",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/sites", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateSite(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.SiteResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "test.com", response.Domain)

	mockSiteUseCase.AssertExpectations(t)
}

func TestSiteHandler_CreateSite_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSiteUseCase := new(MockSiteUseCase)
	handler := handlers.NewSiteHandler(mockSiteUseCase)

	reqBody := map[string]interface{}{}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/sites", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateSite(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
}

func TestKeywordHandler_CreateKeyword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockKeywordUseCase := new(MockKeywordUseCase)
	handler := handlers.NewKeywordHandler(mockKeywordUseCase)

	expectedKeyword := &entities.Keyword{ID: 1, Value: "купить чай", SiteID: 1, GroupID: 1}
	mockKeywordUseCase.On("CreateKeyword", "купить чай", 1, 1).Return(expectedKeyword, nil)

	reqBody := dto.CreateKeywordRequest{
		Value:   "купить чай",
		SiteID:  1,
		GroupID: 1,
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/api/keywords", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateKeyword(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response dto.KeywordResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "купить чай", response.Value)
	assert.Equal(t, 1, response.SiteID)

	mockKeywordUseCase.AssertExpectations(t)
}

func TestKeywordHandler_GetKeywords(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockKeywordUseCase := new(MockKeywordUseCase)
	handler := handlers.NewKeywordHandler(mockKeywordUseCase)

	keywords := []*entities.Keyword{
		{ID: 1, Value: "купить чай", SiteID: 1, GroupID: 1},
		{ID: 2, Value: "купить кофе", SiteID: 1, GroupID: 1},
	}
	mockKeywordUseCase.On("GetKeywordsBySite", 1).Return(keywords, nil)

	req := httptest.NewRequest("GET", "/api/keywords?site_id=1", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetKeywords(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []dto.KeywordResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, "купить чай", response[0].Value)
	assert.Equal(t, "купить кофе", response[1].Value)

	mockKeywordUseCase.AssertExpectations(t)
}

func TestKeywordHandler_GetKeywords_MissingSiteID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockKeywordUseCase := new(MockKeywordUseCase)
	handler := handlers.NewKeywordHandler(mockKeywordUseCase)

	req := httptest.NewRequest("GET", "/api/keywords", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetKeywords(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response dto.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "validation_error", response.Error)
	assert.Contains(t, response.Message, "site_id parameter is required")
}

type MockSiteUseCase struct {
	mock.Mock
}

func (m *MockSiteUseCase) CreateSite(domain string) (*entities.Site, error) {
	args := m.Called(domain)
	return args.Get(0).(*entities.Site), args.Error(1)
}

func (m *MockSiteUseCase) DeleteSite(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSiteUseCase) GetAllSites() ([]*entities.Site, error) {
	args := m.Called()
	return args.Get(0).([]*entities.Site), args.Error(1)
}

func (m *MockSiteUseCase) GetSitesByIDs(ids []int) ([]*entities.Site, error) {
	args := m.Called(ids)
	return args.Get(0).([]*entities.Site), args.Error(1)
}

func (m *MockSiteUseCase) GetKeywordsCount(siteID int) (int, error) {
	args := m.Called(siteID)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockSiteUseCase) GetLastPositionUpdateDate(siteID int) (*time.Time, error) {
	args := m.Called(siteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*time.Time), args.Error(1)
}

type MockKeywordUseCase struct {
	mock.Mock
}

func (m *MockKeywordUseCase) CreateKeyword(value string, siteID int, groupID int) (*entities.Keyword, error) {
	args := m.Called(value, siteID, groupID)
	return args.Get(0).(*entities.Keyword), args.Error(1)
}

func (m *MockKeywordUseCase) DeleteKeyword(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockKeywordUseCase) GetKeywordsBySite(siteID int) ([]*entities.Keyword, error) {
	args := m.Called(siteID)
	return args.Get(0).([]*entities.Keyword), args.Error(1)
}
