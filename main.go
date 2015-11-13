// Command docker-unisync allows syncing of directories to a docker-machine over
// ssh using unison.
package main

import (
	"bytes"
	"errors"
	"github.com/Jeffail/gabs"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
	"strings"
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
		verbose = kingpin.
			Flag("verbose", "Verbose mode.").
			Short('v').
			Bool()
		dest = kingpin.
			Flag("destination", "Destination folder (on the Docker Machine) to sync to.").
			Short('d').
			Default(pwd).
			String()
		source = kingpin.
			Flag("source", "Source folder to sync.").
			Short('s').
			Default(pwd).
			String()
		machineName = kingpin.
				Arg("DOCKER-MACHINE-NAME", "Name of Docker Machine to sync to.").
				Required().
				String()
	)

	// Setup '-h' as an alias for the help flag.
	kingpin.CommandLine.HelpFlag.Short('h')

	kingpin.Parse()

	// Check for `unison`.
	_, unisonPathErr := exec.LookPath("unison")
	if unisonPathErr != nil {
		LogError("could not find `unison`, is it installed?", "See git.io/someurl for information on how to install it.")
	}

	// checkUnisonInstallation(*machineName)

	sshUser := getSshUser(*machineName)
	sshHost := getSshHost(*machineName)
	sshKey := getSshKey(*machineName)

	unisonErr := func() {
		LogError("could not run `unison`, is it installed?", "See git.io/install for information on how to install it.")
	}

	if ok, msg := Sync(sshUser, sshHost, sshKey, *source, *dest, false); ok {
		LogInfo("Watching for changes...")
		Watch(*source, func(id uint64, path string, flags []string) {
			if ok, _ = Sync(sshUser, sshHost, sshKey, *source, *dest, *verbose); ok {
				unisonErr()
			}
		})
	} else {
		LogError("could not run `unison`:\n" + msg)
	}
}

func getSshUser(machine string) string {
	out, err := runDockerMachineCmd("inspect", machine, "--format", "\"{{.Driver.SSHUser}}\"")
	if err != nil {
		LogError(err.Error())
	}

	return strings.Trim(out, " \"\n")
}

func getSshHost(machine string) string {
	out, err := runDockerMachineCmd("ip", machine)
	if err != nil {
		LogError(err.Error())
	}

	return strings.Trim(out, " \"\n")
}

func getSshKey(machine string) string {
	// StorePath has changed in docker-machine 0.5, switch accordingly.
	machineInfo, machineErr := runDockerMachineCmd("inspect", machine)
	if machineErr != nil {
		LogError(machineErr.Error())
	}

	infoJson, infoErr := gabs.ParseJSON([]byte(machineInfo))
	if infoErr != nil {
		LogError("could not load machine info from `docker-machine`.")
	}

	// Docker Machine 0.5 has it in a sub-folder.
	storepath, ok := infoJson.Path("Driver.StorePath").Data().(string)
	if ok {
		return storepath + "/machines/" + infoJson.Path("Driver.MachineName").Data().(string) + "/id_rsa"
	}

	storepath, ok = infoJson.Path("StorePath").Data().(string)

	if !ok {
		LogError("could not determine SSH key path from `docker-machine`.")
	}

	return storepath + "/id_rsa"
}

func checkUnisonInstallation(machine string) {
	_, unisonErr := runDockerMachineCmd("ssh", machine, "which unison")
	if unisonErr != nil {
		// Attempt installation of unison.
		_, installErr := runDockerMachineCmd("ssh", machine, "wget http://www.seas.upenn.edu/~bcpierce/unison/download/releases/unison-2.48.3/unison-2.48.3.tar.gz | tar xz -C /tmp && (cd /tmp/unison-* && make UISTYLE=text && sudo cp unison /usr/local/bin/)")

		if installErr != nil {
			LogError(installErr.Error())
		}

		LogInfo("Installed `unison` in the", machine, "Docker Machine.")
	}
}

func runDockerMachineCmd(cmd ...string) (string, error) {
	var (
		out    bytes.Buffer
		stderr bytes.Buffer
	)

	dm := exec.Command("docker-machine", cmd...)
	dm.Stdout = &out
	dm.Stderr = &stderr
	dmError := dm.Run()
	if dmError != nil {
		return "", errors.New("could not run `docker-machine`: " + stderr.String())
	}
	return out.String(), nil
}

func Sync(user string, host string, key string, source string, dest string, verbose bool) (bool, string) {
	var out bytes.Buffer

	cmd := exec.Command("unison",
		"-auto",
		"-batch",
		"-terse",
		"-watch",
		"-ignore",
		"Name {.git*,*.DS_Store}",
		"-sshargs",
		"-i "+key,
		source,
		"ssh://"+user+"@"+host+"/"+dest)

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
