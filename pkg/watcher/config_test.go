package watcher

import (
	"strings"
	"testing"
)

func TestValidConfig(t *testing.T) {
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
    concurrency: 1
  - apiVersion: extensions/v1beta1
    kind: Deployment
    command: "echo hello"
    concurrency: 4
`

	c, errs := ParseAndValidateConfig(strings.NewReader(conf))

	if errs != nil {
		t.Error(errs)
	}

	if c.Watch[0].ApiVersion != "v1" {
		t.Error("Config wrong:", c)
	}

	if c.Watch[1].ApiVersion != "extensions/v1beta1" {
		t.Error("Config wrong:", c)
	}
}
