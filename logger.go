package main

import (
	"github.com/chinthakagodawita/docker-unisync/Godeps/_workspace/src/github.com/fatih/color"
	"os"
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
	for _, message := range messages {
		color.Cyan(message)
	}
}

func LogDebug(messages ...string) {
	for _, message := range messages {
		color.Cyan(message)
	}
}
