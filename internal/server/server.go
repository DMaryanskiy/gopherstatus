package server

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/DMaryanskiy/gopherstatus/internal/monitor"
	"github.com/DMaryanskiy/gopherstatus/internal/storage"
	"github.com/golang-jwt/jwt/v5"
)

type DBInterface interface {
	CreateService(r *storage.Service) error
	FetchServicesByUser(userID uint) ([]storage.Service, error)
	GetUserByEmail(email string) (storage.User, error)
	LatestResultsByID(limit int, userID uint) ([]storage.CheckResult, error)
	CreateUser(r *storage.User) error
	HistoryForService(limit int, service string) ([]storage.CheckResult, error)
	LatestResults(limit int) ([]storage.CheckResult, error)
}

type Server struct {
	monitor *monitor.Monitor
	mux     *http.ServeMux
	db      DBInterface
}

type contextKey string

const userIDKey contextKey = "user_id"

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			http.Error(w, "Unauthorized: invalid claims", http.StatusUnauthorized)
			return
		}
		userID := uint(claims["user_id"].(float64))
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewServer(m *monitor.Monitor, db *storage.DB) *Server {
	s := &Server{
		monitor: m,
		mux:     http.NewServeMux(),
		db:      db,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	// API handler
	s.mux.Handle("/api/status", s.authMiddleware(http.HandlerFunc(s.handleAPIStatus)))
	s.mux.Handle("/api/history", s.authMiddleware(http.HandlerFunc(s.handleHistoryForService)))
	s.mux.Handle("/api/dashboard", s.authMiddleware(http.HandlerFunc(s.handleDashboardAPI)))
	s.mux.HandleFunc("/api/registration", s.handleRegistration)
	s.mux.Handle("/api/services/create", s.authMiddleware(http.HandlerFunc(s.handleCreateService)))
	s.mux.Handle("/api/services/fetch", s.authMiddleware(http.HandlerFunc(s.handleFetchServicesAPI)))

	// Login logout api
	s.mux.HandleFunc("/api/login", s.handleLogin)
	s.mux.HandleFunc("/api/logout", s.handleLogout)
	s.mux.HandleFunc("/login", s.handleLoginForm)

	// Default form handlers
	s.mux.HandleFunc("/", s.handleIndex)
	s.mux.HandleFunc("/dashboard", s.handleDashboard)
	s.mux.HandleFunc("/register", s.handleRegisterForm)
	s.mux.HandleFunc("/services", s.handleFetchServices)
	s.mux.HandleFunc("/services/create", s.handleServiceCreationForm)
}

func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) Handler() http.Handler {
	return s.mux
}
