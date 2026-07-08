package types

type Config struct {
	Name        string     `yaml:"name"`
	Target      string     `yaml:"target"`
	Duration    string     `yaml:"duration"`
	Concurrency int        `yaml:"concurrency"`
	Scenarios   []Scenario `yaml:"scenarios"`
}

type Scenario struct {
	Name    string  `yaml:"name"`
	Weight  int     `yaml:"weight"`
	Request Request `yaml:"request"`
}

type Request struct {
	Method  string            `yaml:"method"`
	Path    string            `yaml:"path"`
	Headers map[string]string `yaml:"headers"`
	Body    map[string]any    `yaml:"body"`
}
