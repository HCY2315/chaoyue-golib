package mock

import "time"

func NowVOTime() string {
	now := time.Now()
	return now.Format(time.RFC3339)
}
