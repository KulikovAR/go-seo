package entities

import "time"

type Position struct {
	ID        int
	KeywordID int
	SiteID    int
	Rank      int
	URL       string
	Title     string
	Date      time.Time

	Keyword *Keyword
	Site    *Site
}
