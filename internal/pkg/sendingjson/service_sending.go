package sendingjson

import (
	"net/http"
)

type ServiceSend interface {
	Sending(w http.ResponseWriter, r *http.Request, data any) error
}
