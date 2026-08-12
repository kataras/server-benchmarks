//go:build windows

package runner

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// setSysProcAttr is a no-op on Windows: taskkill /T walks the process
// tree by PID, so no process-group setup is needed.
func setSysProcAttr(cmd *exec.Cmd) {}

// signalTree force-kills the process tree. Windows console applications
// have no tree-wide equivalent of SIGTERM, so this is as graceful as it
// gets.
func signalTree(cmd *exec.Cmd) error {
	return taskkill(processPid(cmd))
}

// terminateTree force-kills the process tree.
func terminateTree(cmd *exec.Cmd) error {
	return taskkill(processPid(cmd))
}

// taskkill terminates pid and every descendant.
func taskkill(pid int) error {
	if pid <= 0 {
		return nil
	}

	out, err := exec.Command("taskkill", "/T", "/F", "/PID", strconv.Itoa(pid)).CombinedOutput()
	if err == nil {
		return nil
	}

	// Exit code 128 means the process tree is already gone; message
	// checks cover localized variants only loosely, so the caller treats
	// stop errors as advisory when the process has, in fact, exited.
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 128 {
		return nil
	}

	msg := strings.ToLower(strings.TrimSpace(string(out)))
	if strings.Contains(msg, "not found") {
		return nil
	}

	return fmt.Errorf("taskkill /T /F /PID %d: %w: %s", pid, err, msg)
}
