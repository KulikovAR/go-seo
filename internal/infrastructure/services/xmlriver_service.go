package services

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-seo/internal/domain/entities"
)

type XMLRiverService struct {
	baseURL string
	userID  string
	apiKey  string
	client  *http.Client
}

type SearchRequest struct {
	Query   string
	Page    int
	Device  string
	OS      string
	Ads     bool
	Country string
	Lang    string
	Domain  string
}

type SearchResponse struct {
	XMLName  xml.Name `xml:"yandexsearch"`
	Response Response `xml:"response"`
}

type Response struct {
	Date    string  `xml:"date,attr"`
	Found   int     `xml:"found"`
	Results Results `xml:"results"`
}

type Results struct {
	Grouping Grouping `xml:"grouping"`
	Results  []Result `xml:"result"`
}

type Grouping struct {
	Groups []Group `xml:"group"`
}

type Group struct {
	ID       int   `xml:"id,attr"`
	DocCount int   `xml:"doccount"`
	Docs     []Doc `xml:"doc"`
}

type Doc struct {
	URL         string `xml:"url"`
	Title       string `xml:"title"`
	ContentType string `xml:"contenttype"`
}

type Result struct {
	Position int    `xml:"position"`
	URL      string `xml:"url"`
	Title    string `xml:"title"`
	Type     string `xml:"type"`
}

func NewXMLRiverService(baseURL, userID, apiKey string) (*XMLRiverService, error) {
	return &XMLRiverService{
		baseURL: baseURL,
		userID:  userID,
		apiKey:  apiKey,
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

func (s *XMLRiverService) Search(req SearchRequest, source string) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("user", s.userID)
	params.Set("key", s.apiKey)
	params.Set("query", req.Query)
	
	if source == entities.YandexSearch {
		if req.Page > 1 {
			params.Set("page", strconv.Itoa(req.Page-1))
		}
	} else {
		if req.Page > 0 {
			params.Set("page", strconv.Itoa(req.Page))
		}
	}

	if req.Device != "" {
		params.Set("device", req.Device)
	}
	if req.OS != "" && req.Device == "mobile" {
		params.Set("os", req.OS)
	}
	if req.Ads {
		params.Set("ads", "1")
	}
	if req.Country != "" {
		params.Set("country", req.Country)
	}
	if req.Lang != "" {
		params.Set("lr", req.Lang)
	}
	if req.Domain != "" {
		params.Set("domain", req.Domain)
	}

	var endpoint string
	if source == entities.YandexSearch {
		endpoint = "/search_yandex/xml"
	} else {
		endpoint = "/search/xml"
	}

	requestURL := fmt.Sprintf("%s%s?%s", s.baseURL, endpoint, params.Encode())

	resp, err := s.client.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to XMLRiver: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("XMLRiver API returned status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var searchResp SearchResponse
	if err := xml.Unmarshal(bodyBytes, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w", err)
	}

	return &searchResp, nil
}

func (s *XMLRiverService) findSitePositionInternalWithSubdomains(req SearchRequest, siteDomain string, source string, maxPages int, subdomains bool) (int, string, string, error) {
	for page := 1; page <= maxPages; page++ {
		req.Page = page

		resp, err := s.Search(req, source)
		if err != nil {
			return 0, "", "", fmt.Errorf("failed to search page %d: %w", page, err)
		}

		position := 1
		for _, group := range resp.Response.Results.Grouping.Groups {
			for _, doc := range group.Docs {
				if doc.ContentType == "organic" {
					if s.isSiteMatchWithSubdomains(doc.URL, siteDomain, subdomains) {
						absolutePosition := (page-1)*10 + position
						return absolutePosition, doc.URL, doc.Title, nil
					}
					position++
				}
			}
		}
	}

	return 0, "", "", nil
}
func (s *XMLRiverService) findSitePositionInternal(req SearchRequest, siteDomain string, source string, maxPages int) (int, string, string, error) {
	for page := 1; page <= maxPages; page++ {
		req.Page = page

		resp, err := s.Search(req, source)
		if err != nil {
			return 0, "", "", fmt.Errorf("failed to search page %d: %w", page, err)
		}

		position := 1
		for _, group := range resp.Response.Results.Grouping.Groups {
			for _, doc := range group.Docs {
				if doc.ContentType == "organic" {
					if s.isSiteMatch(doc.URL, siteDomain) {
						absolutePosition := (page-1)*10 + position
						return absolutePosition, doc.URL, doc.Title, nil
					}
					position++
				}
			}
		}
	}

	return 0, "", "", nil
}

func (s *XMLRiverService) FindSitePosition(query, siteDomain, source string, maxPages int, device, os string, ads bool, country, lang string) (int, string, string, error) {
	req := SearchRequest{
		Query:   query,
		Page:    1,
		Device:  device,
		OS:      os,
		Ads:     ads,
		Country: country,
		Lang:    lang,
	}

	return s.findSitePositionInternal(req, siteDomain, source, maxPages)
}

func (s *XMLRiverService) isSiteMatch(resultURL, siteDomain string) bool {
	resultDomain := s.extractDomain(resultURL)
	siteDomainExtracted := s.extractDomain(siteDomain)

	resultDomain = strings.ToLower(strings.TrimPrefix(resultDomain, "www."))
	siteDomainExtracted = strings.ToLower(strings.TrimPrefix(siteDomainExtracted, "www."))

	return resultDomain == siteDomainExtracted
}
func (s *XMLRiverService) FindSitePositionWithSubdomains(query, siteDomain, source string, maxPages int, device, os string, ads bool, country, lang string, subdomains bool) (int, string, string, error) {
	req := SearchRequest{
		Query:   query,
		Page:    1,
		Device:  device,
		OS:      os,
		Ads:     ads,
		Country: country,
		Lang:    lang,
	}

	return s.findSitePositionInternalWithSubdomains(req, siteDomain, source, maxPages, subdomains)
}

func (s *XMLRiverService) isSiteMatchWithSubdomains(resultURL, siteDomain string, subdomains bool) bool {
	resultDomain := s.extractDomain(resultURL)
	siteDomainExtracted := s.extractDomain(siteDomain)

	resultDomain = strings.ToLower(strings.TrimPrefix(resultDomain, "www."))
	siteDomainExtracted = strings.ToLower(strings.TrimPrefix(siteDomainExtracted, "www."))

	exactMatch := resultDomain == siteDomainExtracted
	if exactMatch {
		return true
	}

	if subdomains {
		subdomainMatch := strings.HasSuffix(resultDomain, "."+siteDomainExtracted)
		if subdomainMatch {
			return true
		}

		parentMatch := strings.HasSuffix(siteDomainExtracted, "."+resultDomain)
		if parentMatch {
			return true
		}
	}

	return false
}
func (s *XMLRiverService) extractDomain(urlStr string) string {
	if !strings.HasPrefix(urlStr, "http") {
		urlStr = "http://" + urlStr
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	return strings.ToLower(u.Host)
}

func (s *XMLRiverService) Close() error {
	return nil
}

func (s *XMLRiverService) IsSiteMatch(resultURL, siteDomain string) bool {
	return s.isSiteMatch(resultURL, siteDomain)
}

func (s *XMLRiverService) IsSiteMatchWithSubdomains(resultURL, siteDomain string, subdomains bool) bool {
	return s.isSiteMatchWithSubdomains(resultURL, siteDomain, subdomains)
}
