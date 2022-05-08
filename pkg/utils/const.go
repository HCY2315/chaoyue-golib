package utils

import "fmt"

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

type RunMode string

const (
	DevMode        RunMode = "dev"
	ProductionMode         = "production"
)

var modeMap = map[RunMode]struct{}{
	DevMode:        {},
	ProductionMode: {},
}

func ParseRunMode(rm string) (RunMode, error) {
	mode := RunMode(rm)
	_, find := modeMap[mode]
	if !find {
		return "", fmt.Errorf("未知RunMode")
	}
	return mode, nil
}
