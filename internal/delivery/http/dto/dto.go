package dto

import "time"

type CreateSiteRequest struct {
	Domain string `json:"domain" binding:"required"`
}

type SiteResponse struct {
	ID                 int        `json:"id"`
	Domain             string     `json:"domain"`
	KeywordsCount      int        `json:"keywords_count"`
	LastPositionUpdate *time.Time `json:"last_position_update,omitempty"`
	YandexDynamic      *int       `json:"yandex_dynamic"`
	GoogleDynamic      *int       `json:"google_dynamic"`
}

type DeleteSiteResponse struct {
	Message string `json:"message"`
}

type CreateKeywordRequest struct {
	Value   string `json:"value" binding:"required"`
	SiteID  int    `json:"site_id" binding:"required"`
	GroupID *int   `json:"group_id"`
}

type CreateKeywordItem struct {
	Value   string `json:"value" binding:"required"`
	SiteID  int    `json:"site_id" binding:"required"`
	GroupID *int   `json:"group_id"`
}

type UpdateKeywordRequest struct {
	GroupID *int `json:"group_id"`
}

type KeywordResponse struct {
	ID      int    `json:"id"`
	Value   string `json:"value"`
	SiteID  int    `json:"site_id"`
	GroupID *int   `json:"group_id"`
}

type CreateGroupRequest struct {
	Name   string `json:"name" binding:"required"`
	SiteID int    `json:"site_id" binding:"required"`
}

type UpdateGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type GroupResponse struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	SiteID int    `json:"site_id"`
}

type DeleteKeywordResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

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

type TrackGooglePositionsRequest struct {
	SiteID        int    `json:"site_id" binding:"required"`
	Pages         int    `json:"pages" binding:"omitempty,min=1,max=10"`
	Device        string `json:"device" binding:"omitempty,oneof=desktop tablet mobile"`
	OS            string `json:"os" binding:"omitempty,oneof=ios android"`
	Ads           bool   `json:"ads"`
	Country       string `json:"country"`
	Lang          string `json:"lang"`
	Subdomains    bool   `json:"subdomains"`
	XMLUserID     string `json:"xml_user_id"`
	XMLAPIKey     string `json:"xml_api_key"`
	XMLBaseURL    string `json:"xml_base_url"`
	TBS           string `json:"tbs"`
	Filter        int    `json:"filter"`
	Highlights    int    `json:"highlights"`
	NFPR          int    `json:"nfpr"`
	Loc           int    `json:"loc"`
	AI            int    `json:"ai"`
	Raw           string `json:"raw"`
	LR            int    `json:"lr"`
	Domain        int    `json:"domain"`
	FilterGroupID *int   `json:"filter_group_id"`
}

type TrackYandexPositionsRequest struct {
	SiteID        int    `json:"site_id" binding:"required"`
	Pages         int    `json:"pages" binding:"omitempty,min=1,max=10"`
	Device        string `json:"device" binding:"omitempty,oneof=desktop tablet mobile"`
	OS            string `json:"os" binding:"omitempty,oneof=ios android"`
	Ads           bool   `json:"ads"`
	Country       string `json:"country"`
	Lang          string `json:"lang"`
	Subdomains    bool   `json:"subdomains"`
	XMLUserID     string `json:"xml_user_id"`
	XMLAPIKey     string `json:"xml_api_key"`
	XMLBaseURL    string `json:"xml_base_url"`
	GroupBy       int    `json:"groupby"`
	Filter        int    `json:"filter"`
	Highlights    int    `json:"highlights"`
	Within        int    `json:"within"`
	LR            int    `json:"lr"`
	Raw           string `json:"raw"`
	InIndex       int    `json:"inindex"`
	Strict        int    `json:"strict"`
	Organic       bool   `json:"organic"`
	FilterGroupID *int   `json:"filter_group_id"`
}

type TrackWordstatPositionsRequest struct {
	SiteID                 int    `json:"site_id" binding:"required"`
	XMLUserID              string `json:"xml_user_id"`
	XMLAPIKey              string `json:"xml_api_key"`
	XMLBaseURL             string `json:"xml_base_url"`
	Regions                *int   `json:"regions"`
	Default                *bool  `json:"default"`
	Quotes                 *bool  `json:"quotes"`
	QuotesExclamationMarks *bool  `json:"quotes_exclamation_marks"`
	ExclamationMarks       *bool  `json:"exclamation_marks"`
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

type PositionHistoryResponse struct {
	Data       []PositionHistoryItem `json:"data"`
	Pagination PaginationInfo        `json:"pagination"`
	Meta       MetaInfo              `json:"meta"`
}

type PositionHistoryItem struct {
	ID        int       `json:"id"`
	SiteID    int       `json:"site_id"`
	KeywordID int       `json:"keyword_id"`
	Keyword   string    `json:"keyword"`
	Position  int       `json:"position"`
	Rank      int       `json:"rank"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Date      time.Time `json:"date"`
	Source    string    `json:"source"`
	Device    string    `json:"device"`
	Country   string    `json:"country"`
	Lang      string    `json:"lang"`
}

type PaginationInfo struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	Total       int  `json:"total"`
	LastPage    int  `json:"last_page"`
	From        int  `json:"from"`
	To          int  `json:"to"`
	HasMore     bool `json:"has_more"`
}

type MetaInfo struct {
	QueryTimeMs int  `json:"query_time_ms"`
	Cached      bool `json:"cached"`
}

type PositionHistoryRequest struct {
	SiteID    int     `form:"site_id" binding:"required"`
	KeywordID *int    `form:"keyword_id"`
	Source    *string `form:"source"`
	DateFrom  *string `form:"date_from"`
	DateTo    *string `form:"date_to"`
	Last      *bool   `form:"last"`
	Page      int     `form:"page" binding:"omitempty,min=1"`
	PerPage   int     `form:"per_page" binding:"omitempty,min=1,max=100"`
}

type TrackPositionsResponse struct {
	Message string `json:"message"`
	Count   int    `json:"count"`
}

type AsyncTrackPositionsResponse struct {
	Message string `json:"message"`
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
}

type PositionStatisticsRequest struct {
	SiteID        int    `json:"site_id" binding:"required"`
	DateFrom      string `json:"date_from" binding:"required"`
	DateTo        string `json:"date_to" binding:"required"`
	Source        string `json:"source" binding:"required,oneof=google yandex wordstat"`
	FilterGroupID *int   `json:"filter_group_id"`
}

type PositionStatisticsResponse struct {
	TotalPositions  int             `json:"total_positions"`
	KeywordsCount   int             `json:"keywords_count"`
	Visible         int             `json:"visible"`
	NotVisible      int             `json:"not_visible"`
	PositionRanges  PositionRanges  `json:"position_ranges"`
	VisibilityStats VisibilityStats `json:"visibility_stats"`
	Trends          Trends          `json:"trends"`
}

type PositionRanges struct {
	Range1_3     int `json:"1_3"`
	Range4_10    int `json:"4_10"`
	Range11_30   int `json:"11_30"`
	Range31_50   int `json:"31_50"`
	Range51_100  int `json:"51_100"`
	Range100Plus int `json:"100_plus"`
	NotFound     int `json:"not_found"`
}

type VisibilityStats struct {
	AvgPosition    float64 `json:"avg_position"`
	MedianPosition int     `json:"median_position"`
	BestPosition   int     `json:"best_position"`
	WorstPosition  int     `json:"worst_position"`
}

type Trends struct {
	Improved int `json:"improved"`
	Declined int `json:"declined"`
	Stable   int `json:"stable"`
}

type CombinedPositionsRequest struct {
	SiteID            int     `form:"site_id" binding:"required"`
	Source            *string `form:"source" binding:"omitempty,oneof=google yandex"`
	Wordstat          *bool   `form:"wordstat"`
	WordstatSort      *string `form:"wordstat_sort" binding:"omitempty,oneof=asc desc"`
	DateFrom          *string `form:"date_from"`
	DateTo            *string `form:"date_to"`
	DateSort          *string `form:"date_sort"`
	SortType          *string `form:"sort_type" binding:"omitempty,oneof=asc desc"`
	RankFrom          *int    `form:"rank_from" binding:"omitempty,min=0"`
	RankTo            *int    `form:"rank_to" binding:"omitempty,min=0"`
	Page              int     `form:"page" binding:"omitempty,min=1"`
	PerPage           int     `form:"per_page" binding:"omitempty,min=1,max=100"`
	GroupID           *int    `form:"group_id"`
	FilterGroupID     *int    `form:"filter_group_id"`
	WordstatQueryType *string `form:"wordstat_query_type" binding:"omitempty,oneof=default quotes quotes_exclamation_marks exclamation_marks"`
}

type CombinedPositionsResponse struct {
	Data       []CombinedPositionItem `json:"data"`
	Pagination PaginationInfo         `json:"pagination"`
	Meta       MetaInfo               `json:"meta"`
}

type CombinedPositionItem struct {
	ID        int       `json:"id"`
	SiteID    int       `json:"site_id"`
	KeywordID int       `json:"keyword_id"`
	Keyword   string    `json:"keyword"`
	Date      time.Time `json:"date"`

	Positions []PositionData `json:"positions"`

	Wordstat *PositionData `json:"wordstat"`
}

type TrackingJobsRequest struct {
	SiteID  *int    `form:"site_id"`
	Status  *string `form:"status" binding:"omitempty,oneof=pending running completed failed cancelled"`
	Page    int     `form:"page" binding:"omitempty,min=1"`
	PerPage int     `form:"per_page" binding:"omitempty,min=1,max=100"`
}

type TrackingJobsResponse struct {
	Data       []TrackingJobItem `json:"data"`
	Pagination PaginationInfo    `json:"pagination"`
	Meta       MetaInfo          `json:"meta"`
}

type TrackingJobItem struct {
	ID             string     `json:"id"`
	SiteID         int        `json:"site_id"`
	Source         string     `json:"source"`
	Status         string     `json:"status"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	TotalTasks     int        `json:"total_tasks"`
	CompletedTasks int        `json:"completed_tasks"`
	FailedTasks    int        `json:"failed_tasks"`
	Error          string     `json:"error,omitempty"`
	Progress       float64    `json:"progress"` // Процент выполнения
}

type PositionData struct {
	Rank   int       `json:"rank"`
	URL    string    `json:"url"`
	Title  string    `json:"title"`
	Source string    `json:"source"`
	Date   time.Time `json:"date"`
}
