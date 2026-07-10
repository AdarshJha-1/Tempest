package executor

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/config"
	"github.com/AdarshJha-1/Tempest/internal/metrics"
)

type Executor interface {
	Run(cfg *config.Config) (*metrics.Result, error)
}

type executor struct {
	client *http.Client
}

func New() Executor {
	return &executor{
		client: &http.Client{},
	}
}

func (e *executor) Run(cfg *config.Config) (*metrics.Result, error) {

	wg := sync.WaitGroup{}
	result := &metrics.Result{}

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

func (e *executor) virtualUser(p *testPlan, result *metrics.Result) {

	timer := time.NewTimer(p.duration)

	for {
		select {
		case <-timer.C:
			fmt.Println("work done!")
			return
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
