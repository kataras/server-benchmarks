//go:build !darwin
// +build !darwin

package main

import "os/exec"

func wrapCmd(cmd *exec.Cmd) {}

func killCmd(cmd *exec.Cmd) {
	switch runtime.GOOS {
	case "windows":
		err := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
		if err != nil && err.Error() != "exit status 128" {
			return err
		}

		return nil
	default:
		return cmd.Process.Kill()
		// return exec.Command("kill", "-INT", "-"+strconv.Itoa(cmd.Process.Pid)).Run()
	}
}
