package executor

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/config"
)

type Executor interface {
	Run(cfg *config.Config) (bool, error)
}

type executor struct {
	client *http.Client
}

func New() Executor {
	return &executor{
		client: &http.Client{},
	}
}

func (e *executor) Run(cfg *config.Config) (bool, error) {

	wg := sync.WaitGroup{}

	parsedDuration, err := time.ParseDuration(cfg.Duration)
	if err != nil {
		return false, err
	}
	tp := NewTestPlan(cfg.Target, parsedDuration, cfg.Concurrency, cfg.Scenarios)

	for i := 0; i < tp.concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			e.virtualUser(tp)
		}()
	}
	wg.Wait()

	return true, nil
}

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

func (e *executor) virtualUser(p *testPlan) error {

	timer := time.NewTimer(p.duration)

	for {
		select {
		case <-timer.C:
			fmt.Println("work done!")
			return nil
		default:
			scenario := p.pickScenario()
			if scenario == nil {
				continue
			}
			request, err := p.buildRequest(scenario)
			if err != nil {
				return err
			}
			resp, err := e.client.Do(request)
			time.Sleep(1 * time.Second)
			if err != nil {
				return err
			}
			resp.Body.Close()
		}
	}
}
