package entities

const (
	GoogleSearch = "google"
	YandexSearch = "yandex"
	Wordstat     = "wordstat"
)

type Searcher struct {
	Name string
}
