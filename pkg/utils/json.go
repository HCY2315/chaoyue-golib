package utils

import "strconv"

// 自定义类型上实现 Marshaler 的接口, 在进行 Marshal 时就会使用此除的实现来进行 json 编码
type StrictFloat64 float64

func (f StrictFloat64) MarshalJSON() ([]byte, error) {
	if float64(f) == float64(int(f)) {
		return []byte(strconv.FormatFloat(float64(f), 'f', 2, 64)), nil // 可以自由调整精度
	}
	return []byte(strconv.FormatFloat(float64(f), 'f', -1, 64)), nil
}
