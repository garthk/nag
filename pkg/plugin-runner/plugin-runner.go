package runner

import (
	"bytes"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const DEFAULT_TIMEOUT = 60 * time.Second

type PluginResult struct {
	OriginalOutput string
	Output         string
	Status         PluginStatus
}

type PluginRunOptions struct {
	Timeout   time.Duration
	Sensitive bool
}

func RunPlugin(options PluginRunOptions, command string, args ...string) (PluginResult, error) {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Start(); err != nil {
		return PluginResult{
			OriginalOutput: "",
			Output:         "UNKNOWN\n",
			Status:         UNKNOWN,
		}, err
	}

	// wait or timeout
	donec := make(chan error, 1)
	go func() {
		donec <- cmd.Wait()
	}()

	timeout := options.Timeout
	if timeout.Nanoseconds() == 0 {
		timeout = DEFAULT_TIMEOUT
	}

	select {
	case <-time.After(timeout):
		cmd.Process.Kill()
		return PluginResult{
			OriginalOutput: stdout.String(),
			Output:         "TIMEOUT\n",
			Status:         UNKNOWN,
		}, nil
	case <-donec:
		// pass
	}

	status := cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	pluginStatus := ExitStatusToPluginStatus(status, options.Sensitive)
	output := stdout.String()

	result := PluginResult{
		OriginalOutput: output,
		Output:         DefaultOutputIfBlank(TruncateAndEnforceNewline(output), pluginStatus),
		Status:         pluginStatus,
	}
	return result, nil
}
