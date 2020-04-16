package server

import (
	"net/http"
	"net/http/pprof"
)

//setupRoutes adds http handlers to the servers mux
func (s *Server) setupRoutes() {
	//pprof endpoints for debugging purposes
	s.mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	s.mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	s.mux.Handle("/debug/pps.muxof/profile", http.HandlerFunc(pprof.Profile))
	s.mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	s.mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	s.mux.HandleFunc("/status", s.handleStatus())
	s.mux.HandleFunc("/readings", s.handleReading())
	s.mux.HandleFunc("/stats", s.handleStats())
}

//handleStatus parses the imei query parameter
func (s *Server) handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//finish implementation
	}
}

func (s *Server) handleReading() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//finish implementation
	}
}

func (s *Server) handleStats() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//finish implementation
	}
}
