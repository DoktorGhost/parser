package storage

type Card struct {
	Address              string
	Name                 string
	Category             string
	SubCategory          string
	Url                  string
	UrlImage             string
	Price                string
	PriceWithoutDiscount string
}

type Repo interface {
	Add(name string, card Card) error
}
