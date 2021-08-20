package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/logandavies181/tfd/cmd"
)

var semver string

func main() {
	debugInfo, ok := debug.ReadBuildInfo()
	if semver == "" {
		if !ok {
			fmt.Fprintln(os.Stderr, "Could not determine build info. Your tfd version is corrupt")
			os.Exit(1)
		}

		var buildMeta string
		if debugInfo.Main.Sum == "" {
			buildMeta = "dev"
		} else {
			buildMeta = debugInfo.Main.Sum
		}

		semver = "0.0.0+" + buildMeta
	}

	cmd.Execute(semver)
}
