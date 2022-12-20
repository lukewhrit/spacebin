package util

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
)

// Logger uses zerolog to log information about each request (log level = INFO)
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			log.Info().
				Str("method", r.Method).
				Str("host", r.Host).
				Str("client", r.RemoteAddr).
				Str("page", r.RequestURI).
				Str("protocol", r.Proto).
				Str("user-agent", r.UserAgent()).
				Dur("duration", time.Since(t)).
				Int("status", ww.Status()).
				Int("size", ww.BytesWritten()).
				Msg("HTTP Request")
		}()

		next.ServeHTTP(ww, r)
	})
}
