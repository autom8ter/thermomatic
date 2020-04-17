package server

import (
	"encoding/json"
	"github.com/autom8ter/thermomatic/internal/common"
	"net/http"
	"net/http/pprof"
	"runtime"
	"strconv"
	"time"
)

func (s *Server) setupRoutes() {
	s.mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	s.mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	s.mux.Handle("/debug/pps.muxof/profile", http.HandlerFunc(pprof.Profile))
	s.mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	s.mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	s.mux.HandleFunc("/status", s.handleStatus())
	s.mux.HandleFunc("/readings", s.handleReading())
	s.mux.HandleFunc("/stats", s.handleStats())
}

func (s *Server) handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "expecting method: GET", http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Query().Get("imei")
		if id == "" {
			http.Error(w, "missing imei", http.StatusBadRequest)
			return
		}
		uid, err := strconv.ParseUint(id, 0, 64)
		if err != nil {
			http.Error(w, "invalid uid", http.StatusBadRequest)
			return
		}
		if reading, ok := s.GetReading(uid); ok {
			//if a reading has been stored in the past 5 minutes, return 200
			if time.Since(reading.Timestamp) < 5*time.Minute {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			http.Error(w, "reading not found", http.StatusNotFound)
			return
		}
	}
}

func (s *Server) handleReading() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "expecting method: GET", http.StatusMethodNotAllowed)
			return
		}
		id := r.URL.Query().Get("imei")
		if id == "" {
			http.Error(w, "missing imei", http.StatusBadRequest)
			return
		}
		uid, err := strconv.ParseUint(id, 0, 64)
		if err != nil {
			http.Error(w, "invalid uid", http.StatusBadRequest)
			return
		}
		if reading, ok := s.GetReading(uid); ok {
			if err := json.NewEncoder(w).Encode(reading); err != nil {
				s.serverLog.Printf("failed to encode reading = %s", err.Error())
				http.Error(w, "failed to encode reading", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "reading not found", http.StatusNotFound)
			return
		}
	}
}

//handleStats serves health related metrics.
func (s *Server) handleStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "expecting method: GET", http.StatusMethodNotAllowed)
			return
		}
		stats := &common.Stats{
			GoRoutines:        runtime.NumGoroutine(),
			ClientConnections: s.TotalClients(),
			CPUs:              runtime.NumCPU(),
			Version:           runtime.Version(),
		}

		if err := json.NewEncoder(w).Encode(stats); err != nil {
			s.serverLog.Printf("failed to encode stats = %s", err.Error())
			http.Error(w, "failed to encode stats", http.StatusInternalServerError)
			return
		}
	}
}
