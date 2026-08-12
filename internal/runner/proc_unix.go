//go:build unix

package runner

import (
	"errors"
	"os/exec"
	"syscall"
)

// setSysProcAttr places the child in its own process group, so the whole
// tree — e.g. `go run .` and the compiled server binary it spawns — can be
// signaled together. Without this, killing `go run` alone would leave the
// actual server alive and holding the port.
func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

// signalTree asks the process group to terminate gracefully (SIGTERM).
func signalTree(cmd *exec.Cmd) error {
	return killGroup(cmd, syscall.SIGTERM)
}

// terminateTree kills the process group immediately (SIGKILL).
func terminateTree(cmd *exec.Cmd) error {
	return killGroup(cmd, syscall.SIGKILL)
}

func killGroup(cmd *exec.Cmd, sig syscall.Signal) error {
	pid := processPid(cmd)
	if pid <= 0 {
		return nil
	}

	// Negative PID targets the whole process group (see setSysProcAttr).
	err := syscall.Kill(-pid, sig)
	if errors.Is(err, syscall.ESRCH) {
		return nil // the group is already gone.
	}

	return err
}
