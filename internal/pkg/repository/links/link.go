package links

type Links struct {
	InitialURL string `json:"initial_url"`
	ShortenURL string `json:"shorten_url"`
}

type LinkRepo interface {
	AddLink(item *Links) error
	GetInitialLink(url string) (*Links, error)
	GetShortenLink(url string) (*Links, error)
}
