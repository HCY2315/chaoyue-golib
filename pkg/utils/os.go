package utils

import (
	"os"
	"path"
)

func DefaultHostIP() string {
	hosts, _ := LocalIPv4s()
	return hosts[0]
}

func FileNameInSameDir(targetFilePath string, newFileName string) string {
	dir := path.Dir(targetFilePath)
	return path.Join(dir, newFileName)
}

func HostName(defaultName string) string {
	hn, _ := os.Hostname()
	if hn == "" {
		return defaultName
	}
	return hn
}
