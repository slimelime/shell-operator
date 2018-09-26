package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func SetupShellCommand(ctx context.Context, shellCommand string, envVars map[string]string) *exec.Cmd {
	cmds := strings.Split(shellCommand, " ")
	cmd := exec.CommandContext(ctx, cmds[0], cmds[1:]...)
	cmd.Env = os.Environ()
	for k, v := range envVars {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	return cmd
}

func SetupAndRunShellCommand(ctx context.Context, shellCommand string, envVars map[string]string) ([]byte, error) {
	cmd := SetupShellCommand(ctx, shellCommand, envVars)
	return cmd.CombinedOutput()
}

func RunCommand(cmd *exec.Cmd) error {
	return cmd.Run()
}
