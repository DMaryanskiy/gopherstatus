package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DMaryanskiy/gopherstatus/internal/auth"
	"github.com/DMaryanskiy/gopherstatus/internal/storage"
)

func TestHandleDashboardAPI(t *testing.T) {
	// Mock auth
	origAuth := auth.GetUserIDFromRequest
	auth.GetUserIDFromRequest = func(r *http.Request) (uint, error) {
		return 42, nil
	}
	defer func() { auth.GetUserIDFromRequest = origAuth }()

	s := &Server{db: &MockDB{}}

	req := httptest.NewRequest(http.MethodGet, "/api/dashboard", nil)
	req.Header.Set("Authorization", "Bearer test.token")
	rec := httptest.NewRecorder()

	s.handleDashboardAPI(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var result []storage.CheckResult
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if len(result) != 1 || result[0].ServiceName != "TestService" {
		t.Errorf("unexpected result: %+v", result)
	}
}
