package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"test-task-ozon/internal/pkg/repository/links"
)

const (
	prefixURL = "localhost:8080/"
)

func (h *LinksHandler) Generation(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	link := new(links.Links)
	err := decoder.Decode(&link)
	if err != nil {
		h.Logger.Infof("url:%s method:%s error: failed to decrypt during generation - %s", r.URL.Path, r.Method, err.Error())
		http.Error(w, `generation failed`, http.StatusBadRequest)
		return
	}
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

	res, err := client.Generation(context.Background(), &RequestGeneration{
		InitialUrl: link.InitialURL,
	})
	if err != nil {
		h.Logger.Infof("url:%s method:%s error: generation failed - %s", r.URL.Path, r.Method, err.Error())
		http.Error(w, `generation failed`, http.StatusInternalServerError)
		return
	}
	err = h.Send.Sending(w, r, res.GetShortenUrl())
	if err != nil {
		return
	}
}
