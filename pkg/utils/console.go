package utils

import (
	"fmt"
	"strings"

	"github.com/TwiN/go-color"
)

func WarningMessage(msg string) string {

	return color.Colorize(color.Yellow, msg)
}

func ErrorMessage(msg string) string {

	return color.Colorize(color.Red, msg)
}

var (
	SuccessText = color.Colorize(color.Green, "PASS")
	FailText    = color.Colorize(color.Red, "FAIL")
)

func PrintCuttingLine(msg string) {
	fmt.Printf("%s%s%s\n", strings.Repeat("-", 20), msg, strings.Repeat("-", 20))
}

func PrintSuccess() {
	PrintCuttingLine(SuccessText)
}

func PrintEmptyLine() {
	fmt.Printf("\n")
}

func PrintFail() {
	PrintCuttingLine(FailText)
}

func HumanBytesText(size uint64) string {
	if size >= GB {
		return fmt.Sprintf("%.2f G", float64(size)/float64(GB))
	}
	if size >= MB {
		return fmt.Sprintf("%.2f M", float64(size)/float64(MB))
	}
	if size >= KB {
		return fmt.Sprintf("%.2f K", float64(size)/float64(KB))
	}
	return fmt.Sprintf("%d B", size)
}
