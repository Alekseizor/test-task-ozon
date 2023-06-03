package links

import (
	"database/sql"
	"sync"
)

type RepoLinkInMemory struct {
	links map[string]*Links
	mu    *sync.RWMutex
}

func NewRepoLinkInMemory() (*RepoLinkInMemory, error) {
	return &RepoLinkInMemory{
		links: make(map[string]*Links, 0),
		mu:    &sync.RWMutex{},
	}, nil
}

func (lm *RepoLinkInMemory) AddLink(item *Links) error {
	lm.mu.Lock()
	lm.links[item.ShortenURL] = item
	lm.mu.Unlock()
	return nil
}

func (lm *RepoLinkInMemory) GetInitialLink(url string) (*Links, error) {
	lm.mu.RLock()
	link, existence := lm.links[url]
	lm.mu.RUnlock()
	if existence {
		return link, nil
	}
	return nil, sql.ErrNoRows
}

func (lm *RepoLinkInMemory) GetShortenLink(url string) (*Links, error) {
	lm.mu.RLock()
	for _, initialLink := range lm.links {
		if initialLink.InitialURL == url {
			return initialLink, nil
		}
	}
	lm.mu.RUnlock()
	return nil, nil
}
