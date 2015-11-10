package main

import (
	// "flag"
	"bytes"
	"github.com/chinthakagodawita/docker-unisync/Godeps/_workspace/src/gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
)

var (
	Version string
	Build   string
)

func main() {
	pwd, pwdErr := os.Getwd()
	if pwdErr != nil {
		LogError("could not determine your current directory: " + pwdErr.Error())
	}

	var (
		// verbose     = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
		dest        = kingpin.Flag("destination", "Destination folder (on the Docker Machine) to sync to.").Short('d').Default(pwd).String()
		source      = kingpin.Flag("source", "Source folder to sync.").Short('s').Default(pwd).String()
		machineName = kingpin.Arg("DOCKER-MACHINE-NAME", "Name of Docker Machine to sync to.").Required().String()
	)

	// Setup '-h' as an alias for the help flag.
	kingpin.CommandLine.HelpFlag.Short('h')

	kingpin.Parse()

	// Check for `unison`.
	unisonPath, unisonPathErr := exec.LookPath("unison")
	if unisonPathErr != nil {
		LogError("could not find `unison`, is it installed?", "See git.io/someurl for information on how to install it.")
	}

	unisonCmd := exec.Command("unison", "-batch", "-ignore=Name {.git*,.vagrant/,*.DS_Store}", "-sshargs", "-o StrictHostKeyChecking=no", "-i asd", *source)
	var unisonOut bytes.Buffer
	unisonCmd.Stdout = &unisonOut
	unisonErr := unisonCmd.Run()
	if unisonErr != nil {
		LogError("could not run `unison`: " + unisonErr.Error())
	}

	LogDebug("Out:", unisonOut.String())

	// fmt.Println(verbose)
	LogInfo(*dest)
	LogInfo(*source)
	LogInfo(*machineName)
	LogInfo(unisonPath)
}

func getSshUser(machine string) string {
	var sshUser string

	return sshUser
}

// func runDockerMachineCmd(cmd) exec.Cmd {
// 	dockerMachine := exec.Command("docker-machine", cmd)
// 	var unisonOut bytes.Buffer
// 	unisonCmd.Stdout = &unisonOut
// 	unisonErr := unisonCmd.Run()
// 	if unisonErr != nil {
// 		color.Red("Error: could not run `unison`:")
// 		color.Red(unisonErr.Error())
// 		os.Exit(1)
// 	}
// }
