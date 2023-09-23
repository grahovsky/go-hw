package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/urfave/negroni"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tStart := time.Now()
		lrw := negroni.NewResponseWriter(w)
		next.ServeHTTP(lrw, r)
		status := lrw.Status()

		statusText := http.StatusText(status)
		msg := fmt.Sprintf("%v %v %v %v %v %v %v %v",
			tStart,
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			r.Proto,
			status,
			statusText,
			r.UserAgent(),
		)
		logger.Info(msg)
	})
}
