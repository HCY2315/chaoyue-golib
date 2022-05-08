package errors

import "fmt"

func PanicMsg(recovered interface{}) string {
	var s string
	switch recovered.(type) {
	case string:
		s = recovered.(string)
	case error:
		s = recovered.(error).Error()
	default:
		s = fmt.Sprintf("%+v", recovered)
	}
	return fmt.Sprintf("Panic:%s", s)
}
