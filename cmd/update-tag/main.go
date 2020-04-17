package main

import (
	"flag"
	"fmt"
	"github.com/magbeat/update-tag/services"
	"os"
)

var Version = "development"

func main() {
	var vPrefix = flag.Bool("v", true, "indicates if an 'v' prefix should be added to the tags, default true")
	var prereleasePrefix = flag.String("pre", "RC", "value of the prerelease prefix, default RC")
	var forceDevTags = flag.Bool("forceDev", false, "forces to choose from development tags (ignoring the current branch)")
	var forceProdTags = flag.Bool("forceProd", false, "forces to choose from production tags (ignoring the current branch)")
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
