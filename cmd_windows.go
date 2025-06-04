package main

import (
	"os/exec"
	"syscall"
)

const (
	cmdBinName = "./youtube-dl.exe"
	workingDir = "."
)

func aggrementCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}
