package dto

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
