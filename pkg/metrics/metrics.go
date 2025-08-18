package metrics

import (
	"sync"
	"time"
)

// Simple in-memory HTTP metrics using EWMA for response time and success ratio as uptime proxy
type httpMetrics struct {
	mu              sync.RWMutex
	ewmaRespNs      float64
	ewmaAlpha       float64
	totalRequests   int64
	successRequests int64
}

var defaultHTTP = &httpMetrics{ewmaAlpha: 0.2}

// Observe records one HTTP request result
func ObserveHTTP(status int, duration time.Duration) {
	m := defaultHTTP
	m.mu.Lock()
	// EWMA for response time in nanoseconds
	d := float64(duration.Nanoseconds())
	if m.ewmaRespNs == 0 {
		m.ewmaRespNs = d
	} else {
		m.ewmaRespNs = m.ewmaAlpha*d + (1-m.ewmaAlpha)*m.ewmaRespNs
	}
	m.totalRequests++
	if status >= 200 && status < 500 { // treat 5xx as downtime impact
		m.successRequests++
	}
	m.mu.Unlock()
}

// Snapshot returns avg response time (ms) and uptime percentage based on success ratio
func SnapshotHTTP() (avgMs float64, uptimePerc float64) {
	m := defaultHTTP
	m.mu.RLock()
	ewma := m.ewmaRespNs
	total := m.totalRequests
	success := m.successRequests
	m.mu.RUnlock()
	if ewma > 0 {
		avgMs = ewma / 1e6
	}
	if total > 0 {
		uptimePerc = (float64(success) / float64(total)) * 100
	} else {
		uptimePerc = 100
	}
	return
}
