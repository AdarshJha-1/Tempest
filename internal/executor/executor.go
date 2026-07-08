package executor

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/types"
)

type Executor interface {
	Run(cfg *types.Config) (*types.Result, error)
}

type executor struct {
	client *http.Client
}

func New() Executor {
	return &executor{
		client: &http.Client{},
	}
}

func (e *executor) Run(cfg *types.Config) (*types.Result, error) {

	wg := sync.WaitGroup{}
	result := &types.Result{}

	parsedDuration, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		return nil, err
	}
	tp := NewTestPlan(cfg.Target, parsedDuration, cfg.Concurrency, cfg.Scenarios)

	for i := 0; i < tp.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.virtualUser(tp, result)
		}()
	}
	wg.Wait()

	return result, nil
}

type testPlan struct {
	target      string
	duration    time.Duration
	concurrency int

	scenarios []types.Scenario
}

func NewTestPlan(target string, duration time.Duration, concurrency int, scenarios []types.Scenario) *testPlan {
	return &testPlan{
		target:      target,
		duration:    duration,
		concurrency: concurrency,
		scenarios:   scenarios,
	}
}

func (p *testPlan) buildRequest(s *types.Scenario) (*http.Request, error) {
	var body io.Reader = nil
	if s.Request.Method == "POST" {
	}

	req, err := http.NewRequest(s.Request.Method, p.target+s.Request.Path, body)
	return req, err
}

func (p *testPlan) pickScenario() *types.Scenario {
	randV := rand.Intn(100)

	for i, s := range p.scenarios {
		if randV <= s.Weight {
			return &p.scenarios[i]
		}
	}
	return nil
}

func (e *executor) virtualUser(p *testPlan, result *types.Result) {

	timer := time.NewTimer(p.duration)

	for {
		select {
		case <-timer.C:
			fmt.Println("work done!")
		default:
			scenario := p.pickScenario()
			if scenario == nil {
				continue
			}
			request, err := p.buildRequest(scenario)
			if err != nil {
				continue
			}
			start := time.Now()
			resp, err := e.client.Do(request)
			latency := time.Since(start)
			if err != nil {
				result.RecordNetworkError(latency)
				continue
			}
			result.Record(resp.StatusCode, latency)
			resp.Body.Close()
		}
	}
}
