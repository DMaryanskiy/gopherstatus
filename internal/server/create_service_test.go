package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DMaryanskiy/gopherstatus/internal/auth"
)

func TestHandleCreateService_JSON(t *testing.T) {
	mockDB := &MockDB{}
	server := &Server{db: mockDB}

	// JSON payload
	reqBody := map[string]interface{}{
		"name":     "Test Service",
		"url":      "https://example.com",
		"method":   "GET",
		"interval": 30,
		"body":     "",
		"headers": map[string]string{
			"Content-Type": "application/json",
		},
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/services/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer mock.token")

	// Patch the auth method to return a fixed userID without JWT parsing
	originalAuth := auth.GetUserIDFromRequest
	auth.GetUserIDFromRequest = func(r *http.Request) (uint, error) {
		return 42, nil
	}
	defer func() {
		auth.GetUserIDFromRequest = originalAuth
	}()

	// Execute the handler
	rec := httptest.NewRecorder()
	server.handleCreateService(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", rec.Code)
	}

	if len(mockDB.CreatedService) != 1 {
		t.Fatalf("expected 1 created service, got %d", len(mockDB.CreatedService))
	}

	svc := mockDB.CreatedService[0]
	if svc.Name != "Test Service" || svc.Method != "GET" || svc.URL != "https://example.com" {
		t.Errorf("unexpected service: %+v", svc)
	}
}
