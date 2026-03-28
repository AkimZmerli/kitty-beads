package main

import (
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// newDevProxy creates a reverse proxy handler that forwards requests to the Vite dev server.
// HTTP requests and Vite HMR WebSocket connections are tunneled to viteURL.
// /api/* routes are NOT handled here — they are registered before this handler.
func newDevProxy(viteURL string) http.Handler {
	target, err := url.Parse(viteURL)
	if err != nil {
		log.Fatalf("invalid vite URL %q: %v", viteURL, err)
	}

	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("[dev-proxy] %s %s -> %v (is Vite running on %s?)", r.Method, r.URL.Path, err, viteURL)
		http.Error(w, "Vite dev server unavailable — run: cd frontend && pnpm dev", http.StatusBadGateway)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
			proxyWebSocket(w, r, target)
			return
		}
		r.Host = target.Host
		reverseProxy.ServeHTTP(w, r)
	})
}

// proxyWebSocket tunnels a WebSocket upgrade through to the target (Vite HMR).
func proxyWebSocket(w http.ResponseWriter, r *http.Request, target *url.URL) {
	targetConn, err := net.Dial("tcp", target.Host)
	if err != nil {
		log.Printf("[dev-proxy] ws dial %s: %v", target.Host, err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
		return
	}
	defer targetConn.Close()

	r.Host = target.Host
	if err := r.Write(targetConn); err != nil {
		log.Printf("[dev-proxy] ws write request: %v", err)
		return
	}

	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "WebSocket not supported", http.StatusInternalServerError)
		return
	}
	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		log.Printf("[dev-proxy] hijack: %v", err)
		return
	}
	defer clientConn.Close()

	done := make(chan struct{}, 2)
	go func() { io.Copy(targetConn, clientConn); done <- struct{}{} }()
	go func() { io.Copy(clientConn, targetConn); done <- struct{}{} }()
	<-done
}
