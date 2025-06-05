package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DMaryanskiy/gopherstatus/internal/storage"
)

func TestHandleAPIStatus(t *testing.T) {
	s := &Server{db: &MockDB{}}
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	w := httptest.NewRecorder()

	s.handleAPIStatus(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}

	var results []storage.CheckResult
	if err := json.NewDecoder(w.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if len(results) != 1 || results[0].ServiceName != "TestService" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestHandleHistoryForService_Success(t *testing.T) {
	s := &Server{db: &MockDB{}}
	req := httptest.NewRequest(http.MethodGet, "/api/history?service=WebAPI", nil)
	w := httptest.NewRecorder()

	s.handleHistoryForService(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", w.Code)
	}
}

func TestHandleHistoryForService_MissingParam(t *testing.T) {
	s := &Server{db: &MockDB{}}
	req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
	w := httptest.NewRecorder()

	s.handleHistoryForService(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", w.Code)
	}
}

func TestHandleHistoryForService_DBError(t *testing.T) {
	s := &Server{db: &MockDB{}}
	req := httptest.NewRequest(http.MethodGet, "/api/history?service=fail", nil)
	w := httptest.NewRecorder()

	s.handleHistoryForService(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 Internal Server Error, got %d", w.Code)
	}
}
