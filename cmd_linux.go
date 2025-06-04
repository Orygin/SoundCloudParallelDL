package main

import (
	"os/exec"
	"syscall"
)

const (
	cmdBinName = "yt-dlp"
	workingDir = "/data/"
)

func aggrementCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
