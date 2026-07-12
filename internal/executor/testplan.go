package executor

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/config"
)

type testPlan struct {
	target      url.URL
	duration    time.Duration
	concurrency int

	scenarios []config.Scenario
}

func NewTestPlan(target string, duration time.Duration, concurrency int, scenarios []config.Scenario) *testPlan {

	// i assumed that this is valid cause i did validation before hand
	parsedURL, _ := url.Parse(target)

	return &testPlan{
		target:      *parsedURL,
		duration:    duration,
		concurrency: concurrency,
		scenarios:   scenarios,
	}
}

func (p *testPlan) buildRequest(s *config.Scenario) (*http.Request, error) {

	targetURL := p.target
	targetURL.Path = s.Request.Path
	if s.Request.Query != nil {
		q := targetURL.Query()
		for k, v := range s.Request.Query {
			q.Add(k, v)
		}
		targetURL.RawQuery = q.Encode()

	}
	// idk should i not hardcode this ?
	if s.Request.Method == "GET" {
		// i am strict bi*//
		req, err := http.NewRequest(s.Request.Method, targetURL.String(), nil)
		return req, err
	}

	var body io.Reader

	if s.Request.Body != nil {

		var buff bytes.Buffer
		if err := json.NewEncoder(&buff).Encode(s.Request.Body); err != nil {
			return nil, err
		}

		body = &buff
	}

	req, err := http.NewRequest(s.Request.Method, targetURL.String(), body)

	if s.Request.Headers != nil {
		for k, v := range s.Request.Headers {
			req.Header.Set(k, v)
		}
	}

	return req, err
}

func (p *testPlan) pickScenario() *config.Scenario {
	r := rand.Intn(100)

	current := 0

	for i := range p.scenarios {
		current += p.scenarios[i].Weight

		if r < current {
			return &p.scenarios[i]
		}
	}

	return nil
}
