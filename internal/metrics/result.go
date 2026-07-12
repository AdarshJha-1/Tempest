package metrics

import (
	"sync"
	"time"
)

type Result struct {
	mu sync.Mutex

	TotalRequests int64
	Success2xx    int64
	Client4xx     int64
	Server5xx     int64
	NetworkErrors int64

	TotalLatency time.Duration
}

func (r *Result) Record(status int, latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.TotalRequests++

	switch {
	case status >= 200 && status < 300:
		r.Success2xx++

	case status >= 400 && status < 500:
		r.Client4xx++

	case status >= 500:
		r.Server5xx++
	}

	r.TotalLatency += latency
}

func (r *Result) RecordNetworkError(latency time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.TotalRequests++
	r.NetworkErrors++
	r.TotalLatency += latency
}
