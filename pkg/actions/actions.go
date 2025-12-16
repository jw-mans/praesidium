package actions

import (
	"bytes"
	"os/exec"
	"strings"

	"praesidium/pkg/config"
	"praesidium/pkg/util"
)

// RunOnDisconnectActions executes the configured actions on disconnect event
func RunOnDisconnectActions(actions []config.ActionCfg) {
	for _, action := range actions {
		switch {
		case action.Run != "":
			runCommand(action.Run)
		case action.Log != "":
			util.Error("Action log: %s", action.Log)
		default:
			util.Error("Unknown action configuration: %+v", action)
		}
	}
}

func runCommand(cmd string) {
	util.Error("Running action: %s", cmd)

	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		util.Error("Empty command")
		return
	}

	command := exec.Command(parts[0], parts[1:]...)
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	if err := command.Run(); err != nil {
		util.Error(
			"Action command failed: %v, stderr: %s",
			err,
			stderr.String(),
		)
		return
	}

	if stdout.Len() > 0 {
		util.Info("Action command output: %s", stdout.String())
	}
}
