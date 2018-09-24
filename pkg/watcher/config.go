package watcher

import (
	"io"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/golang/glog"
)

// ShellConfig represents the YAML config file that is used with this app.
type ShellConfig struct {
	Boot []struct {
		Command     string            `yaml:"command"`
		Timeout     int               `yaml:"timeout"`
		Environment map[string]string `yaml:"environment"`
	} `yaml:"boot"`

	Watch []struct {
		ApiVersion  string            `yaml:"apiVersion"`
		Kind        string            `yaml:"kind"`
		Command     string            `yaml:"command"`
		Concurrency int               `yaml:"concurrency"`
		Timeout     int               `yaml:"timeout"`
		Environment map[string]string `yaml:"environment"`
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
	glog.V(7).Infof("Raw config data: %s", data)
	if err != nil {
		return nil, []error{err}
	}

	c := ShellConfig{}
	glog.V(7).Infof("Parsing config...")
	err = yaml.Unmarshal(data, &c)

	if err != nil {
		return nil, []error{err}
	}

	glog.V(7).Infof("Validating config...")
	if errors := validateConfig(&c); len(errors) > 0 {
		return nil, errors
	}

	return &c, nil
}
