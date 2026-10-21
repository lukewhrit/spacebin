/*
 * Copyright 2020-2026 Luke Whritenour

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

package server

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

// ExpireDocuments deletes documents older than the configured expiration age, then again every interval,
// until ctx is cancelled. It returns straight away if expiration is turned off.
func (s *Server) ExpireDocuments(ctx context.Context, interval time.Duration) {
	if s.Config.ExpirationAge <= 0 {
		return
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		s.expireDocuments(ctx)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (s *Server) expireDocuments(ctx context.Context) {
	// Custom documents are always available
	deleted, err := s.Database.DeleteExpiredDocuments(ctx, s.Config.ExpirationAge, s.Config.Documents)

	if err != nil {
		// Shutting down cancels the query, which isn't worth reporting
		if ctx.Err() == nil {
			log.Error().
				Err(err).
				Msg("Error deleting expired documents")
		}

		return
	}

	if deleted > 0 {
		log.Info().
			Int64("deleted", deleted).
			Msg("Deleted expired documents")
	}
}
