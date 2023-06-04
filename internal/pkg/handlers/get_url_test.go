package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"test-task-ozon/internal/pkg/repository/links"
	"test-task-ozon/internal/pkg/sendingjson"
	"testing"
	"time"
)

const (
	POST = "POST"
	GET  = "GET"
)

var (
	client = &http.Client{Timeout: time.Second}
)

type SearchRequest struct {
	//метод запроса
	Method string
	//тело запроса
	Body []byte
	ctx  context.Context
}

type CheckoutResultServer struct {
	Status int
	Data   []byte
}

type TestCase struct {
	Request  SearchRequest
	Response CheckoutResultServer
}
type SearchClient struct {
	// урл внешней системы, куда идти
	URL string
}

func CreateUser() (*LinksHandler, error) {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Println(fmt.Errorf("couldn't create a new logger - %s", err.Error()))
		return nil, err
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	serviceSend := sendingjson.NewServiceSend(logger)
	userHandler := &LinksHandler{
		Logger: logger,
		Send:   serviceSend,
	}
	return userHandler, nil
}

func (srv *SearchClient) CheckoutServer(request SearchRequest) (*CheckoutResultServer, error) {
	searcherReq, err := http.NewRequest(request.Method, srv.URL, bytes.NewBuffer(request.Body)) //nolint:errcheck
	if err != nil {
		return nil, err
	}
	searcherReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(searcherReq)
	if err != nil {
		return nil, fmt.Errorf("unknown error %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := &CheckoutResultServer{}
	result.Status = resp.StatusCode
	result.Data = body
	return result, nil
}

func TestGetLinkOK(t *testing.T) {
	userHandler, err := CreateUser()
	if err != nil {
		log.Println(err)
		return
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	ctx := context.Background()
	repo, err := links.NewRepoLinkPostgres(db, ctx)
	if err != nil {
		log.Println(err)
		return
	}
	userHandler.LinkRepo = repo
	result := []string{"initial_url", "shorten_url"}
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE shorten_url=$1 LIMIT 1;").WithArgs("localhost:8080/141O2_5zsO").WillReturnRows(sqlmock.NewRows(result).AddRow("https://chatbot.theb.ai/#/chat/168573123831355", "localhost:8080/141O2_5zsO"))
	cases := []TestCase{
		{
			Response: CheckoutResultServer{
				Status: http.StatusOK,
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(userHandler.GetLink))
	for caseNum, item := range cases {
		req, err := http.NewRequest(GET, ts.URL, nil)
		if err != nil {
			t.Fatal(err)
		}
		req = mux.SetURLVars(req, map[string]string{
			"URL": "localhost:8080/141O2_5zsO",
		})
		rr := httptest.NewRecorder()
		userHandler.GetLink(rr, req)
		if rr.Code != item.Response.Status {
			t.Errorf("[%d] the status code %d  is different from the expected one %d", caseNum, rr.Code, item.Response.Status)
		}
		if string(rr.Body.Bytes()) != "https://chatbot.theb.ai/#/chat/168573123831355" {
			t.Errorf("[%d] invalid body returned, expected - %s, we have - %s", caseNum, "https://chatbot.theb.ai/#/chat/168573123831355", string(rr.Body.Bytes()))
		}
	}
	ts.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
