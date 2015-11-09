package main

import (
	// "flag"
	"fmt"
	"github.com/fatih/color"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	Version     string
	Build       string
	pwd, pwdErr = os.Getwd()
	verbose     = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	dest        = kingpin.Flag("destination", "Destination folder (on the Docker Machine) to sync to.").Short('d').Required().String()
	source      = kingpin.Flag("source", "Source folder to sync.").Short('s').Default(pwd).String()
	machineName = kingpin.Arg("DOCKER-MACHINE-NAME", "Name of Docker Machine to sync to.").Required().String()
)

func main() {
	if pwdErr != nil {
		color.Set(color.FgRed)
		fmt.Println("Error: could not determine your current directory: ", pwdErr)
		color.Unset()
		os.Exit(1)
	}

	kingpin.CommandLine.HelpFlag.Short('h')

	kingpin.Parse()

	os.Exit(0)

	// app := cli.NewApp()
	// app.Name = "docker-unisync"
	// app.Usage = fmt.Sprintf("Unison-based mounts for your boot2docker docker-machine.")
	// app.Version = fmt.Sprintf("%v-%v", Version, Build)
	// app.ArgsUsage = "DOCKER-MACHINE-NAME..."
	// app.HideHelp = true

	// app.Flags = []cli.Flag{
	// 	cli.BoolFlag{
	// 		Name:  "help, h",
	// 		Usage: "show help",
	// 	},
	// }

	// // app.Action = func(c *cli.Context) {
	// // 	if len(c.Args()) != 1 {
	// // 		color.Red("A Docker Machine name is required, see `docker-unisync help`.")
	// // 		os.Exit(1)
	// // 	}

	// // 	fmt.Println("yooo")
	// // }

	// app.Run(os.Args)

	// pwd, err := os.Getwd()
	// if err != nil {
	// 	color.Set(color.FgRed)
	// 	fmt.Println("Error: could not determine your current directory: ", err)
	// 	color.Unset()
	// 	os.Exit(1)
	// }

	// // Setup usage arguments and options.
	// cmdName := "docker-unisync"
	// flag.Usage = func() {
	// 	fmt.Printf("Usage:\n  %v [options] DOCKER-MACHINE-NAME...", cmdName)
	// 	fmt.Println("\nOptions:")
	// 	flag.PrintDefaults()
	// }

	// help := flag.Bool("help", false, "Show this message")
	// verbose := flag.Bool("verbose", false, "Verbose output")
	// showVersion := flag.Bool("version", false, "Show version")

	// flag.Parse()

	// // Make sure a docker-machine name is specified.
	// if len(flag.Args()) < 1 {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

	// if *showVersion {
	// 	fmt.Printf("%v version %v, build %v\n", cmdName, Version, Build)
	// 	return
	// }

	// if *help {
	// 	flag.Usage()
	// 	return
	// }

	// fmt.Println(*help)
	// fmt.Println(*verbose)
	// fmt.Println("Hello, yo!")
	// fmt.Println(pwd)
}
