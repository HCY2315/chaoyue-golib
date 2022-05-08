package utils

import (
	"time"
)

func NowTs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
