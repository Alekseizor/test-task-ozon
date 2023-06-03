package sending_json

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type SendResponse struct {
	Logger *zap.SugaredLogger
}

func (s *SendResponse) Sending(w http.ResponseWriter, r *http.Request, data any) error {
	postByte, err := json.Marshal(data)
	if err != nil {
		s.Logger.Infof("url:%s method:%s error: failed to Marshal - %s", r.URL.Path, r.Method, err.Error())
		http.Error(w, `failed to provide response data in JSON format`, http.StatusInternalServerError)
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(postByte)
	if err != nil {
		s.Logger.Infof("url:%s method:%s error: failed to write bytes - %s", r.URL.Path, r.Method, err.Error())
		http.Error(w, `failed to write data to HTTP reply`, http.StatusInternalServerError)
		return err
	}
	return nil
}
func NewServiceSend(logger *zap.SugaredLogger) *SendResponse {
	return &SendResponse{Logger: logger}
}
