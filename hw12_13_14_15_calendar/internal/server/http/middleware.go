package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
)

func loggingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tStart := time.Now()
			logger.Info(fmt.Sprintf("%v %v %v %v %v %v %v", tStart, r.RemoteAddr, r.Method, r.URL.Path, r.Proto, r.UserAgent(), middleware.GetReqID(r.Context())))

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				logger.Info(fmt.Sprintf("Request completed. %v %v", ww.Status(), time.Since(tStart)))
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
