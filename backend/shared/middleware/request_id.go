package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// requestIDKey is the context key for the request ID.
type requestIDKey struct{}

// RequestIDHeader is the HTTP header name for the request ID.
const RequestIDHeader = "X-Request-ID"

// RequestID returns a middleware that injects a unique request ID into each request.
// The ID is added to the request context and the response header.
func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if request already has an ID (from upstream proxy)
			id := r.Header.Get(RequestIDHeader)
			if id == "" {
				id = generateID()
			}

			// Add to response header
			w.Header().Set(RequestIDHeader, id)

			// Add to request context
			ctx := context.WithValue(r.Context(), requestIDKey{}, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetRequestID retrieves the request ID from the context.
// Returns empty string if not present.
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey{}).(string); ok {
		return id
	}
	return ""
}

// generateID creates a random 16-character hex string.
func generateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}
