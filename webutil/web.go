package webutil

import (
	"encoding/json"
	"github.com/Direct-Debit/go-commons/errlib"
	"net/http"
)

func ErrorResponder(msg string, code int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(code)
		writer := json.NewEncoder(w)
		errlib.PanicError(writer.Encode(struct {
			Error string `json:"error"`
		}{Error: msg}), "Couldn't encode JSON error")
	}
}
