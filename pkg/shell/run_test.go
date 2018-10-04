package shell_test

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/MYOB-Technology/shell-operator/pkg/shell"
)

func TestNewCommand(t *testing.T) {
	cmd := shell.New(context.Background(), "echo hello world")

	if cmd.Args[0] != "echo" || cmd.Args[1] != "hello" {
		t.Error("incorrect args", cmd.Args)
	}
}

func TestAddingEnvironments(t *testing.T) {
	cmd := exec.Command("echo", "hello")

	if len(cmd.Env) != 0 {
		t.Error("env is not empty to start", cmd.Env)
	}

	envs := map[string]string{
		"ENV1": "A",
		"ENV2": "B",
	}

	shell.AddEnvironment(cmd, envs)

	if len(cmd.Env) != 2 {
		t.Error("env not populated", cmd.Env)
	}

	envs = map[string]string{
		"ENV3": "C",
		"ENV4": "D",
	}

	shell.AddEnvironment(cmd, envs)

	if len(cmd.Env) != 4 {
		t.Error("env not appending", cmd.Env)
	}

	failIfEnvVarMissing(t, "ENV1", "A", cmd.Env)
	failIfEnvVarMissing(t, "ENV2", "B", cmd.Env)
	failIfEnvVarMissing(t, "ENV3", "C", cmd.Env)
	failIfEnvVarMissing(t, "ENV4", "D", cmd.Env)
}

func failIfEnvVarMissing(t *testing.T, key, value string, env []string) {
	found := false
	for _, v := range env {
		if v == fmt.Sprintf("%s=%s", key, value) {
			found = true
		}
	}

	if !found {
		t.Error("missing env var", key, value)
	}
}

func TestRunWithProgress(t *testing.T) {
	cmd := exec.Command("./testdata/test-script")

	sout := bytes.NewBufferString("")
	serr := bytes.NewBufferString("")

	err := shell.RunWithProgress(cmd, sout, serr)

	if err != nil {
		t.Error(err)
	}

	if sout.String() != "first line\nsecond line\n" {
		t.Error("stdout output not correct", sout.String())
	}

	if serr.String() != "error first line\nerror second line\n" {
		t.Error("stderr output not correct", serr.String())
	}
}
