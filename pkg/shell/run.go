package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func New(ctx context.Context, cmd string) *exec.Cmd {
	parts := strings.Split(cmd, " ")
	command := exec.CommandContext(ctx, parts[0], parts[1:]...)
	command.Env = os.Environ()
	return command
}

func AddEnvironment(cmd *exec.Cmd, env map[string]string) *exec.Cmd {
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return cmd
}

func RunWithProgress(cmd *exec.Cmd, stdout io.Writer, stderr io.Writer) error {
	stdoutIn, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrIn, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		io.Copy(stdout, stdoutIn)
	}()

	go func() {
		io.Copy(stderr, stderrIn)
	}()

	return cmd.Wait()
}
