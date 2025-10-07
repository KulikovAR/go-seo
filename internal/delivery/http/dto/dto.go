package dto

import "time"

type CreateSiteRequest struct {
	Name   string `json:"name" binding:"required"`
	Domain string `json:"domain" binding:"required"`
}

type SiteResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Domain string `json:"domain"`
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
	SiteID  int    `json:"site_id" binding:"required"`
	Device  string `json:"device" binding:"required,oneof=desktop tablet mobile"`
	OS      string `json:"os" binding:"omitempty,oneof=ios android"`
	Ads     bool   `json:"ads"`
	Country string `json:"country"`
	Lang    string `json:"lang"`
}

type PositionResponse struct {
	ID        int       `json:"id"`
	KeywordID int       `json:"keyword_id"`
	SiteID    int       `json:"site_id"`
	Rank      int       `json:"rank"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	Keyword   string    `json:"keyword,omitempty"`
	Site      string    `json:"site,omitempty"`
}

type TrackPositionsResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}
