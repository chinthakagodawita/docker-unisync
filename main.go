package main

import (
  "fmt"
  "os"
  "flag"
)

var (
  Version string
  Build string
)

func main() {
  pwd, err := os.Getwd()
  if (err != nil) {
    fmt.Println("Error: could not determine your current directory: ", err)
    os.Exit(1)
  }

  // Setup usage arguments and options.
  cmdName := "docker-unisync"
  flag.Usage = func() {
    fmt.Printf("Example usage:\n\t%v [options] DOCKER-MACHINE-NAME...", cmdName)
    fmt.Println("\nOptions:")
    flag.PrintDefaults()
  }

  help := flag.Bool("help", false, "Show this message")
  verbose := flag.Bool("verbose", false, "Verbose output")
  showVersion := flag.Bool("version", false, "Show version")

  flag.Parse()

  if *showVersion {
    fmt.Printf("%v version %v, build %v\n", cmdName, Version, Build)
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
