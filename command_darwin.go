//go:build darwin
// +build darwin

package main

import (
	"os/exec"
	"syscall"
)

func wrapCmd(cmd *exec.Cmd) {
	// Request the OS to assign process group to the new process, to which all its children will belong
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func killCmd(cmd *exec.Cmd) error {
	// killCommand := exec.Command("kill", "-9", strconv.Itoa(cmd.Process.Pid))
	// killCommand := exec.Command("pkill", "-TERM", "-P", strconv.Itoa(cmd.Process.Pid))
	// return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

	return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
}
