package config

type Boot struct {
	Command     string            `yaml:"command"`
	Timeout     int               `yaml:"timeout"`
	Environment map[string]string `yaml:"environment"`
}

type Watch struct {
	ApiVersion  string            `yaml:"apiVersion"`
	Kind        string            `yaml:"kind"`
	Command     string            `yaml:"command"`
	Concurrency int               `yaml:"concurrency"`
	Timeout     int               `yaml:"timeout"`
	Environment map[string]string `yaml:"environment"`
}

type ShellConfig struct {
	Boot  []Boot  `yaml:"boot,omitempty"`
	Watch []Watch `yaml:"watch,omitempty"`
}
