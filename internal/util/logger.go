/*
 * Copyright 2020-2023 Luke Whritenour

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 *     http://www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
