package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type WordstatService struct {
	baseURL string
	userID  string
	apiKey  string
	client  *http.Client
}

type WordstatRequest struct {
	Query string
}

type WordstatResponse struct {
	Associations []WordstatItem `json:"associations"`
	Popular      []WordstatItem `json:"popular"`
}

type WordstatItem struct {
	IsAssociations bool   `json:"isAssociations"`
	Value          string `json:"value"`
	Text           string `json:"text"`
}

type WordstatPosition struct {
	Query     string
	Frequency int
	Position  int
}

func NewWordstatService(baseURL, userID, apiKey string) (*WordstatService, error) {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     90 * time.Second,
		DisableKeepAlives:   false,
	}

	return &WordstatService{
		baseURL: baseURL,
		userID:  userID,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		},
	}, nil
}

func (s *WordstatService) GetWordstatData(query string, regions *int) (*WordstatResponse, error) {
	params := url.Values{}
	params.Set("user", s.userID)
	params.Set("key", s.apiKey)
	params.Set("query", query)

	if regions != nil {
		params.Set("regions", strconv.Itoa(*regions))
	}

	endpoint := "/wordstat/new/json"
	requestURL := fmt.Sprintf("%s%s?%s", s.baseURL, endpoint, params.Encode())

	resp, err := s.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to Wordstat API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Wordstat API returned status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var wordstatResp WordstatResponse
	if err := json.Unmarshal(bodyBytes, &wordstatResp); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &wordstatResp, nil
}

func (s *WordstatService) GetKeywordFrequency(queryForAPI string, originalQuery string, regions *int) (int, error) {
	resp, err := s.GetWordstatData(queryForAPI, regions)
	if err != nil {
		return 0, err
	}

	for _, item := range resp.Popular {
		if item.Text == originalQuery {
			frequency, err := strconv.Atoi(item.Value)
			if err != nil {
				return 0, fmt.Errorf("failed to parse frequency value: %w", err)
			}
			return frequency, nil
		}
	}

	return 0, nil
}

func (s *WordstatService) GetRelatedKeywords(query string, regions *int) ([]WordstatItem, error) {
	resp, err := s.GetWordstatData(query, regions)
	if err != nil {
		return nil, err
	}

	return resp.Associations, nil
}

func (s *WordstatService) Close() error {
	return nil
}
