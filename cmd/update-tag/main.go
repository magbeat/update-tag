package main

import (
	"flag"
	"fmt"
	"github.com/magbeat/update-tag/services"
	"os"
)

var Version = "development"

func main() {
	var vPrefix = flag.Bool("v", true, "indicates if an 'v' prefix should be added to the tags, defaults to true")
	var prereleasePrefix = flag.String("pre", "RC", "value of the pre-release prefix, default is 'RC', eg. 0.1.0-RC.12")
	var forceDevTags = flag.Bool("forceDev", false, "forces to generate tags for a pre-release (ignoring the current branch), eg. 1.0.0-RC0, 0.1.0-RC.12")
	var forceProdTags = flag.Bool("forceProd", false, "forces to generate tags for a release (ignoring the current branch), eg. 1.0.0, 0.2.0")
	flag.Parse()
	if len(flag.Args()) == 1 && flag.Arg(0) == "version" {
		printVersion()
		os.Exit(0)
	} else if len(flag.Args()) > 0 {
		printUsage()
		os.Exit(0)
	}
	services.UpdateTag(*vPrefix, *prereleasePrefix, *forceProdTags, *forceDevTags)
}

func printVersion() {
	fmt.Println("Version:\t", Version)
}

func printUsage() {
	fmt.Println("\nSubcommands:")
	fmt.Println("  version\tprints the current version")
	fmt.Println("\nFlags:")
	flag.PrintDefaults()
}
