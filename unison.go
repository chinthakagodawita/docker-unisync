package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
)

const SYNC_RETRIES = "3"

func Sync(user string, host string, key string, source string, dest string, ignored []string, verbose bool) (bool, string) {
	var out bytes.Buffer

	ignoredList := ""
	if len(ignored) > 0 {
		ignoredList = strings.Join(ignored, ",")
	}

	cmdArgs := []string{
		"-auto",
		"-batch",
		"-terse",
		"-ignorearchives",
		"-retry",
		SYNC_RETRIES,
		"-prefer",
		source,
		"-sshargs",
		"-i " + key,
	}

	if ignoredList != "" {
		cmdArgs = append(cmdArgs, "-ignore", "Name {"+ignoredList+"}")
	}

	cmdArgs = append(cmdArgs, source, "ssh://"+user+"@"+host+"/"+dest)

	cmd := exec.Command("unison", cmdArgs...)

	// Pipe in stdout/stderr if verbose.
	if verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = &out
		cmd.Stderr = &out
	}

	return cmd.Run() == nil, out.String()
}
