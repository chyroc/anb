package internal

import (
	"github.com/fatih/color"
)

func PrintfGreen(msg string, args ...interface{}) {
	colorGreenIns.Printf(msg, args...)
}

func PrintfWhite(msg string, args ...interface{}) {
	colorWhiteIns.Printf(msg, args...)
}

func PrintfYellow(msg string, args ...interface{}) {
	colorYellowIns.Printf(msg, args...)
}

var (
	colorGreenIns  = color.New(color.FgGreen, color.Bold)
	colorWhiteIns  = color.New(color.FgWhite, color.Bold)
	colorYellowIns = color.New(color.FgYellow, color.Bold)
)
