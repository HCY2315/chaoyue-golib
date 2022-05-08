package utils

import (
	"fmt"
	"syscall"
)

func SysRlimit() string {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err.Error()
	}
	return fmt.Sprintf("[Cur:%d, Max:%d]", rLimit.Cur, rLimit.Max)
}

func SysRlimitEx() (string, uint64, uint64) {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err.Error(), 0, 0
	}
	return fmt.Sprintf("[Cur:%d, Max:%d]", rLimit.Cur, rLimit.Max), rLimit.Cur, rLimit.Max
}
