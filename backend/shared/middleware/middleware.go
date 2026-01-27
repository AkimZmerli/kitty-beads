// Package middleware provides HTTP middleware components for the kitty-beads server.
package middleware

import "net/http"

// Middleware is a function that wraps an http.Handler.
type Middleware func(http.Handler) http.Handler

// Chain creates a new middleware chain from the given middlewares.
// Middlewares are applied in order, so the first middleware wraps the outermost layer.
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
