package main

import (
	"flag"
	"github.com/magbeat/update-tag/services"
)

func main() {
	var vPrefix = flag.Bool("v", true, "indicates if an 'v' prefix should be added to the tags, default true")
	var prereleasePrefix = flag.String("pre", "RC", "value of the prerelease prefix, default RC")
	flag.Parse()
	services.UpdateTag(*vPrefix, *prereleasePrefix)
}
