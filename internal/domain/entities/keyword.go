package entities

type Keyword struct {
	ID      int
	Value   string
	SiteID  int
	GroupID *int

	Site  *Site
	Group *Group
}
