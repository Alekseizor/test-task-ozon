package handlers

import (
	"context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"net/http"
	"test-task-ozon/internal/pkg/repository/links"
	"test-task-ozon/internal/pkg/sendingjson"
)

type LinksHandler struct {
	LinkRepo links.LinkRepo
	Logger   *zap.SugaredLogger
	Send     sendingjson.ServiceSend
}

func (h *LinksHandler) GetLink(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conn, err := grpc.Dial("localhost:9879", opts...)
	if err != nil {
		h.Logger.Infof("url:%s method:%s error: failed to connect to localhost:9879 - %s", r.URL.Path, r.Method, err.Error())
		http.Error(w, `generation failed`, http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	client := NewConverterServiceClient(conn)

	res, err := client.GetLink(context.Background(), &RequestGetLink{
		ShortenUrl: vars["URL"],
	})
	if status.Code(err) == codes.NotFound {
		h.Logger.Infof("url:%s method:%s error: failed to get link - %v", r.URL.Path, r.Method, err)
		http.Error(w, `this link was not found`, http.StatusBadRequest)
		return
	}
	if err != nil {
		h.Logger.Infof("url:%s method:%s error: failed to get link - %v", r.URL.Path, r.Method, err)
		http.Error(w, `couldn't get the original link`, http.StatusInternalServerError)
		return
	}
	err = h.Send.Sending(w, r, res.GetInitialUrl())
	if err != nil {
		return
	}
}
