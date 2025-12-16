package server

import (
	"encoding/json"
	"net/http"

	"praesidium/pkg/monitor"
	"praesidium/pkg/util"
)

type Server struct {
	store *monitor.StatusStore
}

func New(store *monitor.StatusStore) *Server {
	return &Server{
		store: store,
	}
}

func (s *Server) Start(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/status", s.handleStatus)

	util.Info("Starting HTTP server on %s", addr)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil {
			util.Error("HTTP server error: %v", err)
		}
	}()
}

func (s *Server) handleStatus(w http.ResponseWriter, _ *http.Request) {
	status := s.store.Get()
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(status)
}
