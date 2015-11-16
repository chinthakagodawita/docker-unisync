package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"os/exec"
	"strings"
)

func MachineSshUser(machine string) string {
	out, err := RunMachineCommand("inspect", machine, "--format", "\"{{.Driver.SSHUser}}\"")
	if err != nil {
		LogError(err.Error())
	}

	return strings.Trim(out, " \"\n")
}

func MachineIp(machine string) string {
	out, err := RunMachineCommand("ip", machine)
	if err != nil {
		LogError(err.Error())
	}

	host := strings.Trim(out, " \"\n")
	if host == "" {
		LogError(fmt.Sprintf("could not determine IP address of Docker Machine '%v'. Is it running?", machine))
	}

	return host
}

func MachineSshKey(machine string) string {
	// StorePath has changed in docker-machine 0.5, switch accordingly.
	machineInfo, machineErr := RunMachineCommand("inspect", machine)
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

func MachineInstallUnison(machine string) {
	_, unisonErr := RunMachineCommand("ssh", machine, "which unison")
	if unisonErr != nil {
		// Attempt installation of unison.
		_, installErr := RunMachineCommand("ssh", machine, "wget http://www.seas.upenn.edu/~bcpierce/unison/download/releases/unison-2.48.3/unison-2.48.3.tar.gz | tar xz -C /tmp && (cd /tmp/unison-* && make UISTYLE=text && sudo cp unison /usr/local/bin/)")

		if installErr != nil {
			LogError(installErr.Error())
		}

		LogInfo("Installed `unison` in the", machine, "Docker Machine.")
	}
}

func RunMachineCommand(cmd ...string) (string, error) {
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
