package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"test-task-ozon/internal/pkg/repository/links"
	"testing"
)

const (
	shortenURL = "localhost:8080/141O2_5zsO"
	initialURL = "https://chatbot.theb.ai/#/chat/168573123831355"
)

var (
	ctx     = context.Background()
	errTest = fmt.Errorf("test Error")
)

func TestConverterServerGeneration(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	db, mock, err := sqlmock.New()
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	repo, err := links.NewRepoLinkPostgres(db, ctx)
	if err != nil {
		log.Println(err)
		return
	}

	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE initial_url=").WillReturnError(sql.ErrNoRows)
	result := sqlmock.NewResult(1, 1)
	mock.ExpectExec("INSERT INTO link VALUES").WillReturnResult(result)
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE initial_url=").WillReturnError(sql.ErrNoRows)
	mock.ExpectExec("INSERT INTO link VALUES").WillReturnError(errTest)

	svc := ConverterServer{LinkRepo: repo}
	RegisterConverterServiceServer(srv, &svc)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	conn, err := grpc.Dial("", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	t.Cleanup(func() {
		conn.Close()
	})
	if err != nil {
		t.Fatalf("grpc.Dial %v", err)
	}

	client := NewConverterServiceClient(conn)

	res, err := client.Generation(context.Background(), &RequestGeneration{InitialUrl: initialURL})
	if err != nil {
		t.Errorf("client.Generation %v", err)
	}
	if len(res.ShortenUrl) != len(prefixURL)+10 {
		t.Errorf("a different res was expected ShortenUrl %v", res.ShortenUrl)
	}

	res, err = client.Generation(context.Background(), &RequestGeneration{InitialUrl: initialURL})
	if err == nil {
		t.Errorf("another error was expected - %v", errTest)
	}
	if res != nil {
		t.Errorf("expected different result - nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGeneration(t *testing.T) {
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

	repo, err := links.NewRepoLinkPostgres(db, ctx)
	if err != nil {
		log.Println(err)
		return

	}
	go StartConverterServer(repo)

	userHandler.LinkRepo = repo
	result := []string{"initial_url", "shorten_url"}
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE initial_url=").WillReturnRows(sqlmock.NewRows(result).AddRow(initialURL, "141O2_5zsO"))
	cases := []TestCase{
		{
			Request: SearchRequest{
				ctx:    ctx,
				Method: POST,
				Body:   []byte("{\"initial_url\":\"https://chatbot.theb.ai/#/chat/168573123831355\"}"),
			},
			Response: CheckoutResultServer{
				Status: http.StatusOK,
				Data:   []byte("\"" + shortenURL + "\""),
			},
		},
		{
			Request: SearchRequest{
				ctx:    ctx,
				Method: POST,
			},
			Response: CheckoutResultServer{
				Status: http.StatusBadRequest,
				Data:   []byte("generation failed\n"),
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(userHandler.Generation))
	c := &SearchClient{
		URL: ts.URL,
	}
	for caseNum, item := range cases {
		resultServer, err := c.CheckoutServer(item.Request)
		if err != nil {
			log.Println(err)
		}
		if resultServer.Status != item.Response.Status {
			t.Errorf("[%d] the status code %d  is different from the expected one %d", caseNum, resultServer.Status, item.Response.Status)
		}
		if string(resultServer.Data) != string(item.Response.Data) {
			t.Errorf("[%d] invalid body returned, expected - %s, we have - %s", caseNum, string(item.Response.Data), string(resultServer.Data))
		}
	}
	ts.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGenerationInternalError(t *testing.T) {
	userHandler, err := CreateUser()
	if err != nil {
		log.Println(err)
		return
	}
	cases := []TestCase{
		{
			Request: SearchRequest{
				ctx:    ctx,
				Method: POST,
				Body:   []byte("{\"initial_url\":\"https://chatbot.theb.ai/#/chat/168573123831355\"}"),
			},
			Response: CheckoutResultServer{
				Status: http.StatusInternalServerError,
				Data:   []byte("generation failed\n"),
			},
		},
	}
	ts := httptest.NewServer(http.HandlerFunc(userHandler.Generation))
	c := &SearchClient{
		URL: ts.URL,
	}
	for caseNum, item := range cases {
		resultServer, err := c.CheckoutServer(item.Request)
		if err != nil {
			log.Println(err)
		}
		if resultServer.Status != item.Response.Status {
			t.Errorf("[%d] the status code %d  is different from the expected one %d", caseNum, resultServer.Status, item.Response.Status)
		}
		if string(resultServer.Data) != string(item.Response.Data) {
			t.Errorf("[%d] invalid body returned, expected - %s, we have - %s", caseNum, string(item.Response.Data), string(resultServer.Data))
		}
	}
	ts.Close()
}
