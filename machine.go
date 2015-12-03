package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Jeffail/gabs"
	"os/exec"
	"strings"
)

const UNISON_URL = "https://bintray.com/artifact/download/chinthakagodawita/generic/unison-2.48.3-boot2docker.1.tar.gz"

func MachineSshUser(machine string) string {
	machineInfo, machineErr := RunMachineCommand("inspect", machine)
	if machineErr != nil {
		LogError(machineErr.Error())
	}

	infoJson, infoErr := gabs.ParseJSON([]byte(machineInfo))
	if infoErr != nil {
		LogError("could not load machine info from `docker-machine`.")
	}

	user, ok := infoJson.Path("Driver.Driver.SSHUser").Data().(string)
	if ok {
		return user
	}

	user, ok = infoJson.Path("Driver.SSHUser").Data().(string)

	if !ok {
		LogError("could not determine SSH user from `docker-machine`.")
	}

	return user
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

	// Docker Machine 0.5.1 has a direct key.
	storepath, ok := infoJson.Path("Driver.Driver.SSHKeyPath").Data().(string)
	if ok {
		return storepath
	}

	// Docker Machine 0.5 has it in a sub-folder.
	storepath, ok = infoJson.Path("Driver.StorePath").Data().(string)
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
		// Attempt installation of unison (make sure it's a boot2docker machine,
		// error out otherwise).
		machineType, installErr := RunMachineCommand("ssh", machine, "uname -a")
		if installErr != nil {
			LogError(installErr.Error())
		}

		if !strings.Contains(machineType, "boot2docker") {
			LogError("could not find Unison on your Docker Machine.", "Automated Unison installation is not supported on non-boot2docker Docker Machines. Please install Unison on your Docker Machine before continuing.")
		}

		LogInfo("Attempting to install `unison` in your Docker Machine.")

		_, installErr = RunMachineCommand("ssh", machine, "mkdir /tmp/unison && cd /tmp/unison && wget "+UNISON_URL+" && tar xf unison-*.tar.gz && sudo cp unison /usr/local/bin/ && cd /tmp && rm -r /tmp/unison")

		if installErr != nil {
			LogError(installErr.Error())
		}

		LogInfo("Installed `unison` in the", machine, "Docker Machine.")
	}
}

func MachineCreatePath(machine string, path string, user string) {
	_, err := RunMachineCommand("ssh", machine, "mkdir -p", path)
	if err == nil {
		_, err = RunMachineCommand("ssh", machine, "sudo chown -R", user, path)
	}

	if err != nil {
		LogError("could not create destination directory on your Docker Machine: " + err.Error())
	}
}

func RunMachineCommand(cmd ...string) (string, error) {
	var (
		out    bytes.Buffer
		stderr bytes.Buffer
	)

	// LogDebug("running: docker-machine " + strings.Join(cmd, " "))

	dm := exec.Command("docker-machine", cmd...)
	dm.Stdout = &out
	dm.Stderr = &stderr
	dmError := dm.Run()
	if dmError != nil {
		return "", errors.New("could not run `docker-machine`: " + stderr.String())
	}
	return out.String(), nil
}
