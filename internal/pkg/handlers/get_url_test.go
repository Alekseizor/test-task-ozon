package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"test-task-ozon/internal/pkg/repository/links"
	"test-task-ozon/internal/pkg/sendingjson"
)

const (
	POST = "POST"
	GET  = "GET"
)

var (
	client = &http.Client{Timeout: time.Second}
)

type SearchRequest struct {
	Method string
	Body   []byte
	ctx    context.Context
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

func TestConverterServerGetLink(t *testing.T) {
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
	result := []string{"initial_url", "shorten_url"}
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE shorten_url=").WillReturnRows(sqlmock.NewRows(result).AddRow(initialURL, "141O2_5zsO"))
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE shorten_url=").WillReturnError(sql.ErrNoRows)
	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE shorten_url=").WillReturnError(errTest)
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
		t.Fatalf("grpc.DialContext %v", err)
	}

	client := NewConverterServiceClient(conn)

	res, err := client.GetLink(context.Background(), &RequestGetLink{ShortenUrl: shortenURL})
	if err != nil {
		t.Errorf("client.GetLink %v", err)
	}
	if res.InitialUrl != initialURL {
		t.Errorf("a different res was expected InitialUrl %v", res.InitialUrl)
	}

	res, err = client.GetLink(context.Background(), &RequestGetLink{ShortenUrl: shortenURL})
	if status.Code(err) != codes.NotFound {
		t.Errorf("another error code was expected - %v", codes.NotFound)
	}
	if res != nil {
		t.Errorf("expected different result - nil")
	}

	res, err = client.GetLink(context.Background(), &RequestGetLink{ShortenUrl: shortenURL})
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

func TestCreateGRPCClientConnection(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	srv := grpc.NewServer()
	t.Cleanup(func() {
		srv.Stop()
	})

	svc := ConverterServer{}
	RegisterConverterServiceServer(srv, &svc)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("srv.Serve %v", err)
		}
	}()

	_, err := CreateGRPCClientConnection("bufnet")
	if err != nil {
		t.Errorf("another error was expected - %v", "nil")
	}
	_, err = CreateGRPCClientConnection("")
	if err == nil {
		t.Errorf("an error was expected")
	}
}

func TestGetLink(t *testing.T) {
	userHandler, err := CreateUser()
	if err != nil {
		log.Println(err)
		return
	}
	cases := []TestCase{
		{
			Response: CheckoutResultServer{
				Status: http.StatusInternalServerError,
				Data:   []byte("couldn't get the original link\n"),
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
		if rr.Body.String() != string(item.Response.Data) {
			t.Errorf("[%d] invalid body returned, expected - %s, we have - %s", caseNum, string(item.Response.Data), rr.Body.String())
		}
	}
	ts.Close()
}
