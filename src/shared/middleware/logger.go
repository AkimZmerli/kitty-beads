package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// Logger returns a middleware that logs HTTP requests with timing information.
func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			// Process request
			next.ServeHTTP(wrapped, r)

			// Calculate duration
			duration := time.Since(start)

			// Get request ID if available
			requestID := GetRequestID(r.Context())

			// Log the request
			if requestID != "" {
				log.Printf("[%s] %s %s %d %s %d bytes",
					requestID,
					r.Method,
					r.URL.Path,
					wrapped.statusCode,
					duration,
					wrapped.written,
				)
			} else {
				log.Printf("%s %s %d %s %d bytes",
					r.Method,
					r.URL.Path,
					wrapped.statusCode,
					duration,
					wrapped.written,
				)
			}
		})
	}
}
