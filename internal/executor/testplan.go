package executor

import (
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/config"
)

type testPlan struct {
	target      string
	duration    time.Duration
	concurrency int

	scenarios []config.Scenario
}

func NewTestPlan(target string, duration time.Duration, concurrency int, scenarios []config.Scenario) *testPlan {
	return &testPlan{
		target:      target,
		duration:    duration,
		concurrency: concurrency,
		scenarios:   scenarios,
	}
}

func (p *testPlan) buildRequest(s *config.Scenario) (*http.Request, error) {
	var body io.Reader = nil
	if s.Request.Method == "POST" {
	}

	req, err := http.NewRequest(s.Request.Method, p.target+s.Request.Path, body)
	return req, err
}

func (p *testPlan) pickScenario() *config.Scenario {
	randV := rand.Intn(100)

	for i, s := range p.scenarios {
		if randV <= s.Weight {
			return &p.scenarios[i]
		}
	}
	return nil
}
