package util

import (
	"os"
	"syscall"
)

func GetExecDirectory() string {
	file, err := os.Getwd()
	if err == nil {
		return file + "/"
	}

	return ""
}

// CheckProcessExist Will return true if the process with PID exists.
func CheckProcessExist(pid int) bool {
	//查找pid是否存在
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	//往当前进程发送信息0
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false
	}
	return true
}
