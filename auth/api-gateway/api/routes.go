package api

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi/v5/middleware"
)

func (s *Server) routes() {

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Heartbeat("/health"))

	s.router.Handle("/api/v1/user/*", reverseProxy("http://user-service:8080"))

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
