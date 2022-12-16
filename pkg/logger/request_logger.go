package logger

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
)

// New returns a logger middleware for chi, that implements the http.Handler interface.
func RequestLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	if logger == nil {
		return func(next http.Handler) http.Handler { return next }
	}
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// span := trace.SpanFromContext(r.Context())
			// logger = logger.With(
			// 	zap.String("dd.trace_id", convertTraceID(span.SpanContext().TraceID().String())),
			// 	zap.String("dd.span_id", convertTraceID(span.SpanContext().SpanID().String())),
			// 	zap.String("dd.service", "microservice"),
			// 	zap.String("dd.env", "serviceEnv"),
			// 	zap.String("dd.version", version.VERSION),
			// )
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()

			defer func() {
				reqLogger := logger.With(
					zap.String("url", r.URL.String()),
					zap.Int("status_code", ww.Status()),
					zap.String("method", r.Method),
					zap.String("referer", r.Header.Get("Referer")),
					zap.String("request_id", middleware.GetReqID(r.Context())),
					zap.String("useragent", r.Header.Get("User-Agent")),
					zap.String("version", r.Proto),
					zap.Int("bytes", ww.BytesWritten()),
					zap.Int("duration", int(time.Since(t1))),
					zap.Duration("duration_display", time.Since(t1)),
				)
				reqLogger.Info("Served")
			}()
			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}

// Datadog specific
func convertTraceID(id string) string {
	if len(id) < 16 {
		return ""
	}
	if len(id) > 16 {
		id = id[16:]
	}
	intValue, err := strconv.ParseUint(id, 16, 64)
	if err != nil {
		return ""
	}
	return strconv.FormatUint(intValue, 10)
}
