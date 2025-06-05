package monitor

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/DMaryanskiy/gopherstatus/internal/storage"
)

// mockDB mocks the SaveResult method
type mockDB struct {
	saved []*storage.CheckResult
}

func (m *mockDB) SaveResult(result *storage.CheckResult) error {
	m.saved = append(m.saved, result)
	return nil
}

func (m *mockDB) FetchServices() ([]storage.Service, error) {
	return nil, errors.New("not implemented")
}

// roundTripFunc lets us fake an HTTP client response
type roundTripFunc func(req *http.Request) *http.Response

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func TestCheckService_Success(t *testing.T) {
	db := &mockDB{}
	client := &http.Client{
		Transport: roundTripFunc(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(strings.NewReader("OK")),
				Header:     make(http.Header),
			}
		}),
	}

	monitor := &Monitor{
		db:          db,
		http_client: client,
	}

	service := storage.Service{
		Name:     "Test Service",
		Method:   "GET",
		URL:      "http://example.com",
		Interval: 30,
	}

	monitor.checkService(service)

	if len(db.saved) != 1 {
		t.Fatalf("expected 1 saved result, got %d", len(db.saved))
	}

	result := db.saved[0]
	if result.ServiceName != "Test Service" {
		t.Errorf("unexpected service name: %s", result.ServiceName)
	}
	if !result.Online {
		t.Error("expected service to be online")
	}
	if result.ResponseMS < 0 {
		t.Error("expected response time > 0")
	}
}
