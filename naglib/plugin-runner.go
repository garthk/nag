package naglib

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const DEFAULT_TIMEOUT = 60 * time.Second

func RunPlugin(context *PluginContext, command string, args ...string) (PluginResult, error) {
	cfg := context.NagiosConfig
	if cfg == nil {
		cfg = NewNagiosConfig()
	}

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	sysprocattr, err := getSysProcAttr(nil, context)
	if err != nil {
		return pluginError(err)
	}
	if sysprocattr != nil {
		log.Printf("will attempt to run as user %d:%d\n", sysprocattr.Credential.Uid, sysprocattr.Credential.Gid)
		// TODO check privs
		cmd.SysProcAttr = sysprocattr
	}

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

	select {
	case <-time.After(cfg.CommandTimeout):
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
	pluginStatus := ExitStatusToPluginStatus(status, context.ExitStatusTreatment)
	output := stdout.String()

	result := PluginResult{
		OriginalOutput: output,
		Output:         DefaultOutputIfBlank(TruncateAndEnforceNewline(output), pluginStatus),
		Status:         pluginStatus,
	}
	return result, nil
}

func pluginError(err error) (PluginResult, error) {
	return PluginResult{
		OriginalOutput: "",
		Output:         fmt.Sprintf("%s\n", err.Error()),
		Status:         UNKNOWN,
	}, err
}
