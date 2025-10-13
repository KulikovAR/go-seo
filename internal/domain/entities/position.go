package entities

import "time"

type Position struct {
	ID        int
	KeywordID int
	SiteID    int
	Rank      int
	URL       string
	Title     string
	Source    string
	Device    string
	OS        string
	Ads       bool
	Country   string
	Lang      string
	Pages     int
	Date      time.Time

	Keyword *Keyword
	Site    *Site
}
