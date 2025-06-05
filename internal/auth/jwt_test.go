package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	// Set up test secret
	secret := "testsecret"
	os.Setenv("JWT_SECRET", secret)

	// Generate token
	tokenStr, err := GenerateJWT(99)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected a non-empty token string")
	}

	// Parse token
	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Fatalf("failed to parse token: %v", err)
	}
	if !parsedToken.Valid {
		t.Fatal("token is not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("claims are not of type MapClaims")
	}

	// Check user ID
	if uid, ok := claims["user_id"].(float64); !ok || uint(uid) != 99 {
		t.Errorf("expected user_id to be 99, got %v", claims["user_id"])
	}

	// Check expiration
	if exp, ok := claims["exp"].(float64); !ok {
		t.Error("exp claim is not a float64")
	} else if time.Unix(int64(exp), 0).Before(time.Now()) {
		t.Error("token is already expired")
	}
}
