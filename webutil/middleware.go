package webutil

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func DebugRequestMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug(RequestBody(r, false))
		next.ServeHTTP(w, r)
	})
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
