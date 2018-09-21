package executor

import "testing"

func TestParseShellCommandAndEnvs(t *testing.T) {
	cmd := SetupShellCommand("echo hello world", map[string]string{"test1": "a"})

	if cmd.Args[0] != "echo" && cmd.Args[1] != "hello" && cmd.Args[2] != "world" {
		t.Error("Incorrect args:", cmd.Args)
	}
}
