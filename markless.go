package main

import (
	"./markless"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

const (
	version = "0.0.1"
)

var (
	follow   = kingpin.Flag("follow", "Output refreshed as the file changes on disk").Short('f').Bool()
	fileName = kingpin.Arg("file path", "The path of the file to display.").Required().String()
)

func main() {
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Version(version)
	kingpin.Parse()

	status, err := markless.Init(
		markless.WithBuffer(*fileName),
		markless.Follow(*follow),
	).Run()
	if err != nil {
		panic(err)
	}

	os.Exit(status)
	return
}
