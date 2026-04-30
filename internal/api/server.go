package api

import (
	"fmt"
	"net/http"
	"the-unified-document-viewer/internal/config"
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
