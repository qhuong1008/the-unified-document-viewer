package api

import (
	"fmt"
	"net/http"
	"net/url"
	"the-unified-document-viewer/internal/config"

	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	mux    *http.ServeMux
}

func NewServer(cfg *config.Config) *Server {
	mux := http.NewServeMux()

	server := &Server{
		config: cfg,
		mux:    mux,
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("/health", s.healthHandler)
	s.mux.HandleFunc("/", s.homeHandler)
}

// FileProxyHandler creates a Gin handler that proxies file requests to bypass CORS
func FileProxyHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileURL := c.Query("url")
		if fileURL == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "URL parameter is required"})
			return
		}

		// Validate URL
		parsedURL, err := url.Parse(fileURL)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
			return
		}

		// Only allow http and https
		if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only HTTP and HTTPS are allowed"})
			return
		}

		// Make request to external URL
		resp, err := http.Get(fileURL)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to fetch file from external URL"})
			return
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("External server returned status %d", resp.StatusCode)})
			return
		}

		// Copy headers from external response
		contentType := resp.Header.Get("Content-Type")
		if contentType != "" {
			c.Header("Content-Type", contentType)
		}

		// Set CORS headers for the response
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		// Handle OPTIONS preflight
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}

		// Set cache headers
		c.Header("Cache-Control", "public, max-age=3600")

		// Stream the file content to the client
		c.DataFromReader(http.StatusOK, resp.ContentLength, contentType, resp.Body, nil)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"ok"}`)
}

func (s *Server) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "The Unified Document Viewer API")
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	return http.ListenAndServe(addr, s.mux)
}
