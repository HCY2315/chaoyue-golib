//go:build !linux
// +build !linux

package utils

func SysRlimit() string {
	return ""
}

func SysRlimitEx() (string, uint64, uint64) {
	return "", 0, 0
}
