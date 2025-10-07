package entities

type Keyword struct {
	ID     int
	Value  string
	SiteID int

	Site *Site
}
