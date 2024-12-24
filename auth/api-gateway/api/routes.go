package api

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func (s *Server) routes() {

	s.router.Use(middleware.Logger)          // Log API requests
	s.router.Use(middleware.Recoverer)       // Recover from panics without crashing the server
	s.router.Use(middleware.RequestID)       // Add a request ID to the context
	s.router.Use(middleware.RealIP)          // Get the real IP of the request
	s.router.Use(httprate.LimitByIP(100, 1)) // Rate limit requests by IP. It applies to all routes

	// Add a heartbeat route
	s.router.Use(middleware.Heartbeat("/health"))

	s.router.Handle("/api/v1/auth/*", reverseProxy("http://authentication-service"))

}

func reverseProxy(target string) http.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	// Customize the director if needed
	proxy.ModifyResponse = func(resp *http.Response) error {
		// Example: Add custom headers or process the response
		resp.Header.Set("X-Gateway", "Chi API Gateway")
		return nil
	}

	return proxy
}
