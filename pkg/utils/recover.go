package utils

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

var (
	libPackageName string
)

func init() {
	pcs := make([]uintptr, 2)
	_ = runtime.Callers(0, pcs)
	libPackageName = GetPackageName(runtime.FuncForPC(pcs[1]).Name())
}

func GetPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}
	return f
}

func GetCallStacks(err interface{}, pcNum int, skipCallers int) string {
	var buf bytes.Buffer
	pc := make([]uintptr, pcNum)
	max := runtime.Callers(skipCallers, pc)

	for index, pcItem := range pc {
		if index >= max {
			break
		}
		f := runtime.FuncForPC(pcItem)

		fName := f.Name()
		pkgName := GetPackageName(fName)

		if pkgName == libPackageName && strings.HasSuffix(fName, "/utils.RecoverWrapper") {
			continue
		}

		file, line := f.FileLine(pcItem)
		buf.WriteString(fmt.Sprintf("%v\r\n", fName))
		buf.WriteString(fmt.Sprintf("\t%v:%v\r\n", file, line))
	}
	return buf.String()
}
