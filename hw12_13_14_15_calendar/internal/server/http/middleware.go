package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := &wrappedResponseWriter{ResponseWriter: w, status: http.StatusOK}

		tStart := time.Now()
		next.ServeHTTP(writer, r)

		statusText := http.StatusText(writer.status)
		msg := fmt.Sprintf("%v %v %v %v %v %v %v %v", tStart, r.RemoteAddr, r.Method, r.URL.Path, r.Proto, writer.status, statusText, r.UserAgent())
		logger.Info(msg)
	})
}

func (w *wrappedResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
