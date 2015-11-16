package main

import (
	"github.com/fatih/color"
	"os"
	"strings"
)

func LogError(errorMessage string, infoMessages ...string) {
	color.Red("Error: " + errorMessage)

	if infoMessages != nil {
		for _, infoMessage := range infoMessages {
			LogInfo(infoMessage)
		}
	}

	os.Exit(1)
}

func LogInfo(messages ...string) {
	color.Cyan(strings.Join(messages, " "))
}

func LogDebug(messages ...string) {
	color.Yellow(strings.Join(messages, " "))
}
