package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery returns a middleware that recovers from panics and logs the stack trace.
// It returns a 500 Internal Server Error to the client.
func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					// Get request ID for correlation
					requestID := GetRequestID(r.Context())

					// Log the panic with stack trace
					if requestID != "" {
						log.Printf("[%s] PANIC: %v\n%s", requestID, err, debug.Stack())
					} else {
						log.Printf("PANIC: %v\n%s", err, debug.Stack())
					}

					// Return 500 to client
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
