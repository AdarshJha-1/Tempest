package executor

import (
	"io"
	"net/http"

	"github.com/AdarshJha-1/Tempest/internal/config"
	"github.com/AdarshJha-1/Tempest/pkg"
)

type Executor interface {
	Run(cfg *config.Config) (bool, error)
}

type executor struct {
	noOfRoutine int
}

func New(noOfRoutine int) Executor {
	return &executor{noOfRoutine: noOfRoutine}
}

func (e *executor) Run(cfg *config.Config) (bool, error) {
	pkg.PrettyPrintJSON(*cfg)

	for _, s := range cfg.Scenarios {

		req, err := buildRequest(cfg.Target, s)
		if err != nil {
			return false, nil
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return false, nil
		}
		defer resp.Body.Close()
	}

	return true, nil
}

func buildRequest(target string, s config.Scenario) (*http.Request, error) {
	var body io.Reader = nil
	if s.Request.Method == "POST" {
	}

	req, err := http.NewRequest(s.Request.Method, target+s.Request.Path, body)
	return req, err
}
