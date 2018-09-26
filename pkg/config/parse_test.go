package config

import (
	"strings"
	"testing"
)

func TestParsingConfig(t *testing.T) {
	conf := `
boot:
  - command: echo boot1
    timeout: 45
  - command: echo boot2
    environment:
      A: b
      C: d

watch:
  - apiVersion: v1
    kind: Pod
    command: "echo hello"
    timeout: 120
  - apiVersion: extensions/v1beta1
    kind: Deployment
    command: "echo hello"
    concurrency: 4
`

	c, errs := Parse(strings.NewReader(conf))

	if errs != nil {
		t.Error(errs)
	}

	if len(c.Watch) != 2 {
		t.Error("Watches have not been parsed properly", c)
	}

	if c.Watch[0].ApiVersion != "v1" {
		t.Error("Config wrong:", c)
	}

	if c.Watch[1].ApiVersion != "extensions/v1beta1" {
		t.Error("Config wrong:", c)
	}

	if c.Boot[0].Timeout != 45 {
		t.Error("The boot set timeout is not correct", c)
	}

	if c.Boot[1].Timeout != 30 {
		t.Error("The boot default timeout is not correct", c)
	}

	if c.Watch[0].Timeout != 120 {
		t.Error("The watch set timeout is not correct", c)
	}

	if c.Watch[1].Timeout != 1200 {
		t.Error("The watch default timeout is not correct", c)
	}

	if c.Watch[0].Concurrency != 1 {
		t.Error("The watch default concurrency is not correct", c)
	}

	if c.Watch[1].Concurrency != 4 {
		t.Error("The watch set concurrency is not correct", c)
	}
}

func TestParsingBadConfig(t *testing.T) {
	conf := `
bojh d

sd fgsdfg
`

	_, err := Parse(strings.NewReader(conf))

	if err == nil {
		t.Error("Should error on parsing.")
	}
}

func TestInvalidConfig(t *testing.T) {
	confs := [][]string{
		[]string{"", "must have at least one watch configured"},
		[]string{"boot:\n- command: test", "must have at least one watch configured"},
		[]string{"watch:\n- timeout: 123", "must have command set for watch config (0)"},
		[]string{"watch:\n- kind: Pod", "must have command set for watch config (0)"},
		[]string{"watch:\n- apiVersion: v1\n  command: echo hello", "must have kind set for watch config (0)"},
		[]string{"watch:\n- kind: Pod\n  command: echo hello", "must have API Version set for watch config (0)"},
		[]string{"watch:\n- apiVersion: v1\n  kind: Pod", "must have command set for watch config (0)"},
		[]string{"watch: []", "must have at least one watch configured"},
	}

	for _, conf := range confs {
		c, err := Parse(strings.NewReader(conf[0]))

		if err == nil {
			t.Error("Invalid config does not provoke error", conf, c)
		}

		if err.Error() != conf[1] {
			t.Error("invalid error message", err.Error(), conf[1])
		}
	}
}
