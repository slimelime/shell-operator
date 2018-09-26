package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func ParseFromFile(path string) (*ShellConfig, error) {
	confIn, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	return Parse(confIn)
}

func Parse(in io.Reader) (*ShellConfig, error) {
	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}

	c := ShellConfig{}
	err = yaml.Unmarshal(data, &c)

	if err != nil {
		return nil, err
	}

	if len(c.Watch) < 1 {
		return nil, errors.New("must have at least one watch configured")
	}

	for i, b := range c.Boot {
		if b.Command == "" {
			return nil, fmt.Errorf("must have command set for boot config (%d)", i)
		}

		if b.Timeout == 0 {
			c.Boot[i].Timeout = 30
		}
	}

	for i, w := range c.Watch {
		if w.Command == "" {
			return nil, fmt.Errorf("must have command set for watch config (%d)", i)
		}

		if w.Kind == "" {
			return nil, fmt.Errorf("must have kind set for watch config (%d)", i)
		}

		if w.ApiVersion == "" {
			return nil, fmt.Errorf("must have API Version set for watch config (%d)", i)
		}

		if w.Timeout == 0 {
			c.Watch[i].Timeout = 1200
		}

		if w.Concurrency == 0 {
			c.Watch[i].Concurrency = 1
		}
	}

	return &c, nil
}
