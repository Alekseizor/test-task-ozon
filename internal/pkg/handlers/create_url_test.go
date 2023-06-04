package handlers

//
//import (
//	"github.com/DATA-DOG/go-sqlmock"
//	"github.com/gorilla/mux"
//	"log"
//	"net/http"
//	"net/http/httptest"
//	"test-task-ogit zon/internal/pkg/repository/links"
//	"testing"
//)
//
//func TestGenerationOK(t *testing.T) {
//	userHandler, err := CreateUser()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	defer db.Close()
//	ctx := context.Background()
//	repo, err := links.NewRepoLinkPostgres(db, ctx)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	userHandler.LinkRepo = repo
//	result := []string{"initial_url", "shorten_url"}
//	mock.ExpectQuery("SELECT initial_url,shorten_url FROM link WHERE shorten_url=$1 LIMIT 1;").WithArgs("localhost:8080/141O2_5zsO").WillReturnRows(sqlmock.NewRows(result).AddRow("https://chatbot.theb.ai/#/chat/168573123831355", "localhost:8080/141O2_5zsO"))
//	cases := []TestCase{
//		{
//			Response: CheckoutResultServer{
//				Status: http.StatusOK,
//			},
//		},
//	}
//	ts := httptest.NewServer(http.HandlerFunc(userHandler.GetLink))
//	for caseNum, item := range cases {
//		req, err := http.NewRequest(POST, ts.URL, nil)
//		if err != nil {
//			t.Fatal(err)
//		}
//		req = mux.SetURLVars(req, map[string]string{
//			"URL": "localhost:8080/141O2_5zsO",
//		})
//		rr := httptest.NewRecorder()
//		userHandler.GetLink(rr, req)
//		if rr.Code != item.Response.Status {
//			t.Errorf("[%d] the status code %d  is different from the expected one %d", caseNum, rr.Code, item.Response.Status)
//		}
//		if string(rr.Body.Bytes()) != "https://chatbot.theb.ai/#/chat/168573123831355" {
//			t.Errorf("[%d] invalid body returned, expected - %s, we have - %s", caseNum, "https://chatbot.theb.ai/#/chat/168573123831355", string(rr.Body.Bytes()))
//		}
//	}
//	ts.Close()
//
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("there were unfulfilled expectations: %s", err)
//	}
//}
