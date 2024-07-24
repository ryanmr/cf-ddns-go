package pkg

import (
	"net/http"
	"sync/atomic"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

var requestID int64

// RequestLogger is a custom middleware that logs requests using Zerolog
//
// It provides a scoped logger for each request/response lifecycle,
// along with an incrementing request id (vs traditional uuid).
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		currentRequestID := atomic.AddInt64(&requestID, 1)

		logger := log.With().Int64("request-id", currentRequestID).Str("method", req.Method).Str("url", req.URL.String()).Logger()

		writer := middleware.NewWrapResponseWriter(res, req.ProtoMajor)

		next.ServeHTTP(writer, req)

		logger.Info().
			Str("module", "http").
			Int("status", writer.Status()).
			Int("bytes", writer.BytesWritten()).
			Int64("duration", time.Since(start).Milliseconds()).Send()

	})
}
