package links

import (
	"database/sql"
	"log"
	"testing"
)

func TestNewRepoLinkInMemory(t *testing.T) {
	_, err := NewRepoLinkInMemory()
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
		return
	}
}

func TestAddLinkInMemory(t *testing.T) {
	repo, err := NewRepoLinkInMemory()
	if err != nil {
		log.Println(err)
		return
	}
	link := &Links{
		InitialURL: initialURL,
		ShortenURL: shortenURL,
	}
	err = repo.AddLink(link)
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
	}
}

func TestGetInitialLinkInMemory(t *testing.T) {
	repo, err := NewRepoLinkInMemory()
	if err != nil {
		log.Println(err)
		return
	}
	repo.links[shortenURL] = &Links{
		ShortenURL: shortenURL,
		InitialURL: initialURL,
	}
	cases := []TestRepoLink{
		{
			url: shortenURL,
			response: TestRepoLinkResponse{
				link: &Links{
					InitialURL: initialURL,
					ShortenURL: shortenURL,
				},
				err: nil,
			},
		},
		{
			url: initialURL,
			response: TestRepoLinkResponse{
				link: nil,
				err:  sql.ErrNoRows,
			},
		},
	}
	for number, testCase := range cases {
		link, err := repo.GetInitialLink(testCase.url)
		if err != testCase.response.err {
			log.Println(err)
			t.Errorf("[%d] the error is different from the expected one", number)
			continue
		}
		if link == nil && testCase.response.link == nil {
			continue
		}
		if (link == nil && testCase.response.link != nil) || (link != nil && testCase.response.link == nil) {
			log.Println(err)
			t.Errorf("[%d] the link is different from the expected one", number)
			continue
		}
		if link.InitialURL != testCase.response.link.InitialURL {
			t.Errorf("[%d] the InitialURL is different from the expected one", number)
		}
		if link.ShortenURL != testCase.response.link.ShortenURL {
			t.Errorf("[%d] the ShortenURL is different from the expected one", number)
		}
	}
}

func TestGetShortenLinkInMemory(t *testing.T) {
	repo, err := NewRepoLinkInMemory()
	if err != nil {
		log.Println(err)
		return
	}
	repo.links[shortenURL] = &Links{
		ShortenURL: shortenURL,
		InitialURL: initialURL,
	}
	cases := []TestRepoLink{
		{
			url: initialURL,
			response: TestRepoLinkResponse{
				link: &Links{
					InitialURL: initialURL,
					ShortenURL: shortenURL,
				},
				err: nil,
			},
		},
		{
			url: shortenURL,
			response: TestRepoLinkResponse{
				link: nil,
				err:  nil,
			},
		},
	}
	for number, testCase := range cases {
		link, err := repo.GetShortenLink(testCase.url)
		if err != testCase.response.err {
			log.Println(err)
			t.Errorf("[%d] the error is different from the expected one", number)
			continue
		}
		if link == nil && testCase.response.link == nil {
			continue
		}
		if (link == nil && testCase.response.link != nil) || (link != nil && testCase.response.link == nil) {
			log.Println(err)
			t.Errorf("[%d] the link is different from the expected one", number)
			continue
		}
		if link.InitialURL != testCase.response.link.InitialURL {
			t.Errorf("[%d] the InitialURL is different from the expected one", number)
		}
		if link.ShortenURL != testCase.response.link.ShortenURL {
			t.Errorf("[%d] the ShortenURL is different from the expected one", number)
		}
	}
}
