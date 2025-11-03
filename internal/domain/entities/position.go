package entities

import "time"

type Position struct {
	ID            int
	KeywordID     int
	SiteID        int
	Rank          int
	URL           string
	Title         string
	Source        string
	Device        string
	OS            string
	Ads           bool
	Country       string
	Lang          string
	Pages         int
	Date          time.Time
	FilterGroupID *int

	Keyword *Keyword
	Site    *Site
}

type PositionStatistics struct {
	TotalPositions       int
	KeywordsCount        int
	Visible              int
	NotVisible           int
	PositionDistribution PositionDistribution
	PositionRanges       PositionRanges
	VisibilityStats      VisibilityStats
	Trends               Trends
}

type PositionDistribution struct {
	Top3     int
	Top10    int
	Top20    int
	NotFound int
}

type PositionRanges struct {
	Range1_3     int
	Range4_10    int
	Range11_30   int
	Range31_50   int
	Range51_100  int
	Range100Plus int
	NotFound     int
}

type VisibilityStats struct {
	AvgPosition    float64
	MedianPosition int
	BestPosition   int
	WorstPosition  int
}

type Trends struct {
	Improved int
	Declined int
	Stable   int
}

type CombinedPosition struct {
	ID        int
	SiteID    int
	KeywordID int
	Keyword   *Keyword
	Date      time.Time

	Positions []*Position

	Wordstat *Position
}
