package main

import (
  "fmt"
  "os"
  "flag"
)

var version string

func main() {
  pwd, err := os.Getwd()
  if (err != nil) {
    fmt.Println("Error: could not determine your current directory: ", err)
    os.Exit(1)
  }

  // Setup usage arguments and options.
  cmdName := "docker-unisync"
  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Example usage:\n\t%v [options] DOCKER-MACHINE-NAME...\n", cmdName)
    fmt.Print("\nOptions:\n")
    flag.PrintDefaults()
  }

  help := flag.Bool("help", false, "Show this message")
  verbose := flag.Bool("verbose", false, "Verbose output")
  showVersion := flag.Bool("version", false, "Show version")

  flag.Parse()

  if *showVersion {
    fmt.Println(cmdName, version)
    return
  }

  if *help {
    flag.Usage()
    return
  }

  fmt.Println(*help)
  fmt.Println(*verbose)
  fmt.Println("Hello, yo!")
  fmt.Println(pwd)
}
