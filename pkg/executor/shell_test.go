package executor

import (
	"context"
	"testing"
	"time"
)

func TestParseShellCommandAndEnvs(t *testing.T) {
	ctx := context.Background()

	cmd := SetupShellCommand(ctx, "echo hello world", map[string]string{"test1": "a"})

	if cmd.Args[0] != "echo" && cmd.Args[1] != "hello" && cmd.Args[2] != "world" {
		t.Error("Incorrect args:", cmd.Args)
	}
}

func TestRunCommandWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	cmd := SetupShellCommand(ctx, "sleep 1", map[string]string{})

	_, err := cmd.Output()

	if err == nil {
		t.Error("Expect command error.")
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Error("No context deadline when expected.", ctx.Err())
	}
}

func TestSetupAndRunShellCommand(t *testing.T) {
	ctx := context.Background()

	output, err := SetupAndRunShellCommand(ctx, "echo hello world", map[string]string{"test1": "a"})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if string(output[:]) != "hello world\n" {
		t.Error("unexpected output", string(output[:]))
	}
}
