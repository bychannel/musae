package elog

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	procSetStdHandle = kernel32.MustFindProc("SetStdHandle")
)

// RedirectStderr to the file passed in
func RedirectStderr(service string) (err error) {
	logFile, err := os.OpenFile(fmt.Sprintf("./%s_error.log", service), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	var r1 uintptr
	var e1 syscall.Errno
	if runtime.GOOS == "windows" {
		eHandle := syscall.STD_ERROR_HANDLE
		r1, _, e1 = syscall.SyscallN(procSetStdHandle.Addr(), 2, uintptr(eHandle), uintptr(syscall.Handle(logFile.Fd())), 0)
	}
	if r1 == 0 {
		if e1 != 0 {
			return error(e1)
		}
		return error(syscall.EINVAL)
	}
	os.Stderr = logFile
	return
}
