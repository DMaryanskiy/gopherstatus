package auth

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(secret string, userID uint) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	})
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestGetUserIDFromRequest_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	token := generateTestToken("testsecret", 42)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	userID, err := GetUserIDFromRequest(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if userID != 42 {
		t.Errorf("expected user ID 42, got %d", userID)
	}
}

func TestGetUserIDFromRequest_MissingHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	_, err := GetUserIDFromRequest(req)
	if err == nil || err.Error() != "missing auth header" {
		t.Errorf("expected missing auth header error, got %v", err)
	}
}

func TestGetUserIDFromRequest_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not.a.valid.token")

	_, err := GetUserIDFromRequest(req)
	if err == nil || err.Error() != "invalid token" {
		t.Errorf("expected invalid token error, got %v", err)
	}
}
