package config

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type ValidationError struct {
	Errors []string
}

func (v *ValidationError) add(msg string) {
	v.Errors = append(v.Errors, msg)
}
func (v *ValidationError) Error() string {
	var b strings.Builder

	fmt.Fprintln(&b, "\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(&b, "  Tempest Configuration Validation Failed")
	fmt.Fprintln(&b, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(&b)

	for i, err := range v.Errors {
		fmt.Fprintf(&b, " %2d. %s\n", i+1, err)
	}

	fmt.Fprintln(&b)
	fmt.Fprintf(&b, "Found %d validation error(s).\n", len(v.Errors))
	fmt.Fprintln(&b, "Fix the configuration and try again.")
	fmt.Fprintln(&b, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	return b.String()
}

func Validate(cfg *Config) error {
	var errs ValidationError

	validateGeneral(cfg, &errs)
	validateTarget(cfg, &errs)
	validateDuration(cfg, &errs)
	validateConcurrency(cfg, &errs)
	validateScenarios(cfg, &errs)
	if len(errs.Errors) > 0 {
		return &errs
	}

	return nil
}

func validateGeneral(cfg *Config, errs *ValidationError) {
	if cfg.Name == "" {
		errs.add("Name is empty")
	}
}
func validateTarget(cfg *Config, errs *ValidationError) {
	if cfg.Target == "" {
		errs.add("Target is empty")
	}
	if _, err := url.ParseRequestURI(cfg.Target); err != nil {
		errs.add(fmt.Sprintf("Invalid target, %s", err))
	}
}
func validateDuration(cfg *Config, errs *ValidationError) {
	if cfg.Duration == "" {
		errs.add("Duration is empty")
	}
	if _, err := time.ParseDuration(cfg.Duration); err != nil {
		errs.add(fmt.Sprintf("Invalid duration, %s", err))
	}
}
func validateConcurrency(cfg *Config, errs *ValidationError) {
	if cfg.Concurrency == 0 || cfg.Concurrency > 100 {
		errs.add("Concurrency must be in range of 1 - 100 inclusive")
	}
}

var validMethods = map[string]bool{
	"GET":    true,
	"POST":   true,
	"PUT":    true,
	"PATCH":  true,
	"DELETE": true,
}

func validateScenarios(cfg *Config, errs *ValidationError) {
	if len(cfg.Scenarios) == 0 {
		errs.add("At least one scenario is required")
		return
	}

	totalWeight := 0

	for i, s := range cfg.Scenarios {
		if s.Name == "" {
			errs.add(fmt.Sprintf("scenario[%d]: name can't be empty", i+1))
		}
		if s.Weight == 0 || s.Weight > 100 {
			errs.add(fmt.Sprintf("scenario[%d]: weight must be in range of 1 - 100 inclusive", i+1))
		}
		method := strings.ToUpper(s.Request.Method)
		if _, ok := validMethods[method]; !ok {
			errs.add(fmt.Sprintf("scenario[%d]: invalid method name", i+1))
		}

		if s.Request.Path != "" && !strings.HasPrefix(s.Request.Path, "/") {
			errs.add(fmt.Sprintf("scenario[%d]: path must start with '/' ", i+1))
		}

		totalWeight += s.Weight
	}

	if totalWeight != 100 {
		errs.add(fmt.Sprintf("scenario weights must sum to 100, got %d", totalWeight))
	}
}
