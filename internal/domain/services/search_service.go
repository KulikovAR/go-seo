package services

type SearchService interface {
	FindSitePosition(query, siteDomain, source string, maxPages int, device, os string, ads bool, country, lang string) (int, string, string, error)
	Close() error
}
