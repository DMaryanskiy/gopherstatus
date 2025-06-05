package monitor

import "time"

type ServiceStatus struct {
	Name       string
	URL        string
	Method     string
	Online     bool
	ResponseMS int64
	CheckedAt  time.Time
	Error      string
}
