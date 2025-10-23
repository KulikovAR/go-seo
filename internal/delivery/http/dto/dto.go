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

// Google-specific DTO
type TrackGooglePositionsRequest struct {
	SiteID     int    `json:"site_id" binding:"required"`
	Pages      int    `json:"pages" binding:"omitempty,min=1,max=10"`
	Device     string `json:"device" binding:"omitempty,oneof=desktop tablet mobile"`
	OS         string `json:"os" binding:"omitempty,oneof=ios android"`
	Ads        bool   `json:"ads"`
	Country    string `json:"country"`
	Lang       string `json:"lang"`
	Subdomains bool   `json:"subdomains"`
	XMLUserID  string `json:"xml_user_id"`
	XMLAPIKey  string `json:"xml_api_key"`
	XMLBaseURL string `json:"xml_base_url"`
	// Google-specific parameters
	TBS        string `json:"tbs"`        // Период поиска: qdr:h, qdr:d, qdr:w, qdr:m, qdr:y
	Filter     int    `json:"filter"`     // Скрывать похожие результаты: 0 или 1
	Highlights int    `json:"highlights"` // Подсветка ключевых слов: 0 или 1
	NFPR       int    `json:"nfpr"`       // Отменить исправление запроса: 0 или 1
	Loc        int    `json:"loc"`        // ID местоположения
	AI         int    `json:"ai"`         // Парсинг блока "Обзор от ИИ": 0 или 1
	Raw        string `json:"raw"`        // Полный HTML код страницы: "page"
}

// Yandex-specific DTO
type TrackYandexPositionsRequest struct {
	SiteID     int    `json:"site_id" binding:"required"`
	Pages      int    `json:"pages" binding:"omitempty,min=1,max=10"`
	Device     string `json:"device" binding:"omitempty,oneof=desktop tablet mobile"`
	OS         string `json:"os" binding:"omitempty,oneof=ios android"`
	Ads        bool   `json:"ads"`
	Country    string `json:"country"`
	Lang       string `json:"lang"`
	Subdomains bool   `json:"subdomains"`
	XMLUserID  string `json:"xml_user_id"`
	XMLAPIKey  string `json:"xml_api_key"`
	XMLBaseURL string `json:"xml_base_url"`
	// Yandex-specific parameters
	GroupBy    int    `json:"groupby"`    // ТОП позиций для сбора (всегда 10)
	Filter     int    `json:"filter"`     // Скрывать похожие результаты: 0 или 1
	Highlights int    `json:"highlights"` // Подсветка ключевых слов: 0 или 1
	Within     int    `json:"within"`     // Фильтр по периоду: 77, 1, 2, 0
	LR         int    `json:"lr"`         // ID региона Яндекса
	Raw        string `json:"raw"`        // Полный HTML код страницы: "page"
	InIndex    int    `json:"inindex"`    // Проверка индексации: 0 или 1
	Strict     int    `json:"strict"`     // Режим строгого соответствия: 0 или 1
}

// Wordstat-specific DTO
type TrackWordstatPositionsRequest struct {
	SiteID     int    `json:"site_id" binding:"required"`
	XMLUserID  string `json:"xml_user_id"`
	XMLAPIKey  string `json:"xml_api_key"`
	XMLBaseURL string `json:"xml_base_url"`
	Regions    *int   `json:"regions"` // ID региона Яндекса, nullable
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

// Async tracking response with task ID
type AsyncTrackPositionsResponse struct {
	Message string `json:"message"`
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
}
