package links

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
)

type TestRepoLinkResponse struct {
	link *Links
	err  error
}

type TestRepoLink struct {
	url      string
	response TestRepoLinkResponse
}

const (
	shortenURL = "localhost:8080/141O2_5zsO"
	initialURL = "https://chatbot.theb.ai/#/chat/168573123831355"
)

var (
	ctx       = context.Background()
	testError = fmt.Errorf("test Error")
)

func TestNewRepoLinkPostgres(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	_, err = NewRepoLinkPostgres(db, ctx)
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
		return
	}
}

func TestAddLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	repo, err := NewRepoLinkPostgres(db, ctx)
	if err != nil {
		log.Println(err)
		return
	}
	link := &Links{
		InitialURL: initialURL,
		ShortenURL: shortenURL,
	}
	result := sqlmock.NewResult(1, 1) // вставляем одну запись, затронуто одна строка
	mock.ExpectExec("INSERT INTO link VALUES").WillReturnResult(result)
	mock.ExpectExec("INSERT INTO link VALUES").WillReturnError(testError)

	err = repo.AddLink(link)
	if err != nil {
		t.Errorf("[0] the error is different from the expected one %s", "nil")
	}

	err = repo.AddLink(link)
	if err != testError {
		t.Errorf("[1] the error is different from the expected one %v", testError)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetInitialLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	repo, err := NewRepoLinkPostgres(db, ctx)
	if err != nil {
		log.Println(err)
		return
	}
	result := []string{"initial_url", "shorten_url"}
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE shorten_url=").WillReturnRows(sqlmock.NewRows(result).AddRow(initialURL, shortenURL))
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE shorten_url=").WillReturnError(testError)
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
			url: shortenURL,
			response: TestRepoLinkResponse{
				link: nil,
				err:  testError,
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
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetShortenLink(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	repo, err := NewRepoLinkPostgres(db, ctx)
	if err != nil {
		log.Println(err)
		return
	}
	result := []string{"initial_url", "shorten_url"}
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE initial_url=").WillReturnRows(sqlmock.NewRows(result).AddRow(initialURL, shortenURL))
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE initial_url=").WillReturnError(testError)
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE initial_url=").WillReturnError(sql.ErrNoRows)
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
			url: initialURL,
			response: TestRepoLinkResponse{
				link: nil,
				err:  testError,
			},
		},
		{
			url: initialURL,
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
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
