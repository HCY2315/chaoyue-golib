package common

import (
	"fmt"
	"strconv"
	"strings"
)

func UintListFromStr(s string, sep string) ([]uint, error) {
	trimmedStr := strings.TrimSpace(s)
	if trimmedStr == "" {
		return nil, fmt.Errorf("(%s) empty", s)
	}
	ss := strings.Split(trimmedStr, ",")
	ret := make([]uint, 0, len(ss))
	for _, idStr := range ss {
		id64, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("(%s) not uint", idStr)
		}
		ret = append(ret, uint(id64))
	}
	return ret, nil
}
