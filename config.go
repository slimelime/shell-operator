package main

import (
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// ShellConfig represents the YAML config file that is used with this app.
type ShellConfig struct {
	Boot struct {
		Command     string            `yaml:"command"`
		Environment map[string]string `yaml:"environment`
	} `yaml:"boot,omitempty"`

	Watch []struct {
		ApiVersion  string            `yaml:"apiVersion"`
		Kind        string            `yaml:"kind"`
		Command     string            `yaml:"command"`
		Concurrency int               `yaml:"concurrency"`
		Environment map[string]string `yaml:"environment`
	} `yaml:"watch"`
}

// TODO: proper validation beyond YAML syntax
func validateConfig(c *ShellConfig) []error {
	return []error{}
}

// ParseAndValidateConfig will read in and marshell the YAML config for running the
// shell operator and either return a valid ready to go config struct, or provide a
// list of validation errors that can be rendered to the user to help them fix their
// config.
func ParseAndValidateConfig(in io.Reader) (*ShellConfig, []error) {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, []error{err}
	}

	c := ShellConfig{}
	err = yaml.Unmarshal(data, &c)

	if err != nil {
		return nil, []error{err}
	}

	if errors := validateConfig(&c); len(errors) > 0 {
		return nil, errors
	}

	return &c, nil
}
