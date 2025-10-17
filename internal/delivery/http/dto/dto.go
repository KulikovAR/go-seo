package dto

import "time"

type CreateSiteRequest struct {
	Domain string `json:"domain" binding:"required"`
}

type SiteResponse struct {
	ID            int    `json:"id"`
	Domain        string `json:"domain"`
	KeywordsCount int    `json:"keywords_count"`
}

type DeleteSiteResponse struct {
	Message string `json:"message"`
}

type CreateKeywordRequest struct {
	Value  string `json:"value" binding:"required"`
	SiteID int    `json:"site_id" binding:"required"`
}

type KeywordResponse struct {
	ID     int    `json:"id"`
	Value  string `json:"value"`
	SiteID int    `json:"site_id"`
}

type DeleteKeywordResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// Position DTOs
type TrackSitePositionsRequest struct {
	SiteID     int    `json:"site_id" binding:"required"`
	Source     string `json:"source" binding:"required,oneof=google yandex wordstat"`
	Pages      int    `json:"pages" binding:"omitempty,min=1,max=10"`
	Device     string `json:"device" binding:"omitempty,oneof=desktop tablet mobile"`
	OS         string `json:"os" binding:"omitempty,oneof=ios android"`
	Ads        bool   `json:"ads"`
	Country    string `json:"country"`
	Lang       string `json:"lang"`
	Subdomains bool   `json:"subdomains"`
}

type PositionResponse struct {
	ID        int       `json:"id"`
	KeywordID int       `json:"keyword_id"`
	SiteID    int       `json:"site_id"`
	Rank      int       `json:"rank"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Source    string    `json:"source"`
	Device    string    `json:"device"`
	OS        string    `json:"os"`
	Ads       bool      `json:"ads"`
	Country   string    `json:"country"`
	Lang      string    `json:"lang"`
	Pages     int       `json:"pages"`
	Date      time.Time `json:"date"`
	Keyword   string    `json:"keyword,omitempty"`
	Site      string    `json:"site,omitempty"`
}

type TrackPositionsResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}
