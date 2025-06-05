package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DMaryanskiy/gopherstatus/internal/auth"
	"github.com/DMaryanskiy/gopherstatus/internal/storage"
)

func TestHandleFetchServicesAPI(t *testing.T) {
	// Mock auth to always return userID = 1
	orig := auth.GetUserIDFromRequest
	auth.GetUserIDFromRequest = func(r *http.Request) (uint, error) {
		return 1, nil
	}
	defer func() {
		auth.GetUserIDFromRequest = orig
	}()

	server := &Server{db: &MockDB{}}

	req := httptest.NewRequest(http.MethodGet, "/api/services/fetch", nil)
	req.Header.Set("Authorization", "Bearer test.token") // simulate token

	rec := httptest.NewRecorder()
	server.handleFetchServicesAPI(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", rec.Code)
	}

	var services []storage.Service
	if err := json.NewDecoder(rec.Body).Decode(&services); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if len(services) != 1 || services[0].Name != "MockService" {
		t.Errorf("unexpected response: %+v", services)
	}
}
