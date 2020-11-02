package webutil

import (
	"encoding/json"
	"fmt"
	"github.com/Direct-Debit/go-commons/errlib"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
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

func PerformanceLogMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(r.URL)
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Info(fmt.Sprintf("Handled in %v seconds", duration.Seconds()))
	})
}
