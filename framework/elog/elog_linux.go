//go:build linux

package elog

import (
	"fmt"
	"os"
	"syscall"
)

// RedirectStderr to the file passed in
func RedirectStderr(service string) (err error) {
	logFile, err := os.OpenFile(fmt.Sprintf("./%s_error.log", service), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	if err != nil {
		return
	}

	err = syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		return
	}

	return
}
