package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cubic-cwnd/internal/check"
	"cubic-cwnd/internal/cubic"
	"cubic-cwnd/internal/fair"
	"cubic-cwnd/internal/input"
	"cubic-cwnd/internal/runbook"
	"cubic-cwnd/internal/sim"
)

type Server struct {
	mux  *http.ServeMux
	addr string
	book *runbook.Book
}

func New(addr string) *Server {
	s := &Server{
		mux:  http.NewServeMux(),
		addr: addr,
		book: runbook.NewBook(64),
	}
	s.routes()
	return s
}

func Serve(addr string) error {
	return New(addr).ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Book() *runbook.Book {
	return s.book
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/win", s.handleWin)
	s.mux.HandleFunc("/api/curve", s.handleCurve)
	s.mux.HandleFunc("/api/sim", s.handleSim)
	s.mux.HandleFunc("/api/fair", s.handleFair)
	s.mux.HandleFunc("/api/check", s.handleCheck)
	s.mux.HandleFunc("/api/history", s.handleHistory)
	s.mux.HandleFunc("/api/health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET only")
		return
	}
	writeJSON(w, http.StatusOK, s.book.List())
}

func (s *Server) handleWin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	spec, err := input.LoadBytes(body)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	res, err := cubic.ComputeFastConv(spec.ToParams(), spec.PrevWMax, true)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	entry := runbook.Entry{
		ID:   fmt.Sprintf("run-%d", s.book.NextSeq()+1),
		Spec: *spec,
		Win:  res,
		Note: spec.Name,
	}
	if err := s.book.Add(entry); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res = cubic.HoldWinLive(res)
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleCurve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	spec, err := input.LoadBytes(body)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	n := spec.Samples
	if n < 2 {
		n = 2
	}
	step := spec.HorizonSeconds / float64(n-1)
	pts := make([]runbook.CurvePoint, 0, n)
	for i := 0; i < n; i++ {
		p := spec.ToParams()
		p.T = step * float64(i)
		pts = append(pts, runbook.CurvePoint{
			T:        p.T,
			WCubic:   cubic.WCubic(p),
			WEst:     cubic.WEst(p),
			W:        cubic.WEffective(p),
			Friendly: cubic.IsTCPFriendly(p),
		})
	}
	writeJSON(w, http.StatusOK, pts)
}

func (s *Server) handleSim(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	spec, err := input.LoadBytes(body)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	mode, err := sim.ParseMode(spec.Sim.Mode)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	start, err := sim.ParseStart(spec.Sim.Start)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	cfg := sim.Config{
		Mode:        mode,
		Start:       start,
		Rounds:      spec.Sim.Rounds,
		InitialCwnd: spec.Sim.InitialCwnd,
		Ssthresh:    spec.Sim.Ssthresh,
		WMax:        spec.WMax,
		C:           spec.C,
		RTT:         spec.RTT,
	}
	res, err := sim.Run(cfg)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleFair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	spec, err := input.LoadBytes(body)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	cfg := fair.Config{
		Capacity: spec.Fair.Capacity,
		Rounds:   spec.Fair.Rounds,
		FlowA:    spec.Fair.FlowA,
		FlowB:    spec.Fair.FlowB,
	}
	res, err := fair.Run(cfg)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, res)
}

func (s *Server) handleCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	spec, err := input.LoadBytes(body)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	results, pass := check.Verify(spec)
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"pass":    pass,
		"results": results,
	})
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
