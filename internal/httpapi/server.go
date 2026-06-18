package httpapi

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/keepdevops/cofiswarm-kvpool/internal/policy"
)

type Server struct {
	mu  sync.RWMutex
	cfg policy.Config
}

func New(cfg policy.Config) *Server {
	return &Server{cfg: cfg}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/v1/policy", func(w http.ResponseWriter, _ *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		_ = json.NewEncoder(w).Encode(s.cfg)
	})
	mux.HandleFunc("/v1/evaluate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var req policy.EvaluateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.mu.RLock()
		defer s.mu.RUnlock()
		_ = json.NewEncoder(w).Encode(s.cfg.Evaluate(req))
	})
	mux.HandleFunc("/v1/admit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var req policy.AdmitRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		s.mu.RLock()
		defer s.mu.RUnlock()
		_ = json.NewEncoder(w).Encode(s.cfg.Admit(req))
	})
	return mux
}
