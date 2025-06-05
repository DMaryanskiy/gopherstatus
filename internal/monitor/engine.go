package monitor

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/DMaryanskiy/gopherstatus/internal/storage"
)

type DBInterface interface {
	SaveResult(r *storage.CheckResult) error
	FetchServices() ([]storage.Service, error)
}

type Monitor struct {
	http_client *http.Client
	db          DBInterface
}

func NewMonitor(db *storage.DB) *Monitor {
	return &Monitor{
		db:          db,
		http_client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *Monitor) Start() {
	services, err := m.db.FetchServices()
	if err != nil {
		log.Fatalf("Failed to load services: %v", err)
	}
	for _, svc := range services {
		go m.monitorService(svc)
	}
}

func (m *Monitor) monitorService(svc storage.Service) {
	ticker := time.NewTicker(time.Duration(svc.Interval) * time.Second)
	defer ticker.Stop()

	for {
		m.checkService(svc)
		<-ticker.C
	}
}

func (m *Monitor) checkService(svc storage.Service) {
	start := time.Now()

	var reqBody io.Reader
	if svc.Body != "" {
		reqBody = strings.NewReader(svc.Body)
	}

	req, err := http.NewRequest(svc.Method, svc.URL, reqBody)
	if err != nil {
		log.Printf("failed to create request for %s: %v", svc.Name, err)
		m.db.SaveResult(&storage.CheckResult{
			ServiceName: svc.Name,
			Method:      svc.Method,
			URL:         svc.URL,
			Online:      false,
			ResponseMS:  0,
			CheckedAt:   time.Now(),
			Error:       err.Error(),
		})
		return
	}

	for _, header := range svc.Headers {
		req.Header.Set(header.Key, header.Value)
	}

	resp, err := m.http_client.Do(req)
	if err != nil {
		log.Printf("failed to check %s: %v", svc.Name, err)
		m.db.SaveResult(&storage.CheckResult{
			ServiceName: svc.Name,
			Method:      svc.Method,
			URL:         svc.URL,
			Online:      false,
			ResponseMS:  0,
			CheckedAt:   time.Now(),
			Error:       err.Error(),
		})
		return
	}

	duration := time.Since(start).Milliseconds()
	online := resp.StatusCode >= 200 && resp.StatusCode < 300

	result := &storage.CheckResult{
		ServiceName: svc.Name,
		Method:      svc.Method,
		URL:         svc.URL,
		Online:      online,
		ResponseMS:  duration,
		CheckedAt:   time.Now(),
		Error:       "",
	}

	m.db.SaveResult(result)
}
