// Command docker-unisync allows syncing of directories to a docker-machine over
// ssh using unison.
package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/exec"
	"strings"
)

var (
	Name    string
	Version string
	Build   string
)

const IGNORE_WILDCARD = "*"

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
		noWatch = kingpin.
			Flag("no-watch", "Don't watch source directory for changes (do a onetime sync).").
			Short('w').
			Default("false").
			Bool()
		ignored = kingpin.
			Flag("ignored", "Comma-separated list of file patterns to ignore (use '"+IGNORE_WILDCARD+"' as a wildcard).").
			Short('i').
			Default("*.git*,.DS_Store").
			String()
		machineName = kingpin.
				Arg("DOCKER-MACHINE-NAME", "Name of Docker Machine to sync to.").
				Required().
				String()
	)

	// Setup '-h' as an alias for the help flag.
	kingpin.CommandLine.HelpFlag.Short('h')

	// Setup version printing.
	kingpin.Flag("version", "Show version.").PreAction(kingpin.Action(func(*kingpin.ParseContext) error {
		fmt.Printf("%v version %v, build %v\n", Name, Version, Build)
		os.Exit(0)
		return nil
	})).Bool()

	kingpin.Parse()

	// Check for `unison`.
	_, unisonPathErr := exec.LookPath("unison")
	if unisonPathErr != nil {
		LogError("could not find `unison`, is it installed?", "See git.io/someurl for information on how to install it.")
	}

	// Check ignored flags.
	var ignoredItems []string
	if *ignored != "" {
		ignoredItems = strings.Split(*ignored, ",")
		for index, ignoredItem := range ignoredItems {
			ignoredItems[index] = strings.Trim(ignoredItem, " ")
		}
	}

	// checkUnisonInstallation(*machineName)

	sshUser := MachineSshUser(*machineName)
	sshHost := MachineIp(*machineName)
	sshKey := MachineSshKey(*machineName)

	unisonErr := func() {
		LogError("could not run `unison`, is it installed?", "See git.io/install for information on how to install it.")
	}

	LogInfo("Beginning initial sync, please wait...")
	if ok, msg := Sync(sshUser, sshHost, sshKey, *source, *dest, ignoredItems, false); ok {
		if !*noWatch {
			LogInfo("Initial sync complete, watching for changes...")
			syncing := false
			Watch(*source, ignoredItems, func(id uint64, path string, flags []string) {
				if !syncing {
					syncing = true
					if ok, _ = Sync(sshUser, sshHost, sshKey, *source, *dest, ignoredItems, *verbose); !ok {
						unisonErr()
					}
				}
				syncing = false
			})
		} else {
			LogInfo("Initial sync complete")
		}
	} else {
		LogError("could not run `unison`:\n" + msg)
	}
}
