package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cubic-cwnd/internal/cubic"
)

func TestWinEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"w_max":16,"rtt_seconds":0.1,"t_seconds":0.4}`))
	req := httptest.NewRequest(http.MethodPost, "/api/win", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	var out cubic.Result
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out.W <= 0 || out.K <= 0 {
		t.Fatalf("win=%+v", out)
	}
	if s.Book().Len() != 1 {
		t.Fatalf("book len=%d", s.Book().Len())
	}
}

func TestCurveEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"w_max":16,"rtt_seconds":0.1,"t_seconds":0.4,"horizon_seconds":1,"samples":5}`))
	req := httptest.NewRequest(http.MethodPost, "/api/curve", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
	var out []map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if len(out) != 5 {
		t.Fatalf("curve len=%d", len(out))
	}
}

func TestSimEndpoint(t *testing.T) {
	s := New(":0")
	body := bytes.NewReader([]byte(`{"w_max":16,"rtt_seconds":0.1,"t_seconds":0.4,"sim":{"mode":"cubic","rounds":10,"start":"after-loss"}}`))
	req := httptest.NewRequest(http.MethodPost, "/api/sim", body)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "States") {
		t.Fatalf("sim body=%s", rec.Body.String())
	}
}

func TestHealthEndpoint(t *testing.T) {
	s := New(":0")
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()
	s.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
}
