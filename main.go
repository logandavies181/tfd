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

		debugVersion := debugInfo.Main.Version
		if debugVersion == "(devel)" {
			debugVersion = "devel"
		}

		semver = "0.0.0+"+debugVersion
	}

	cmd.Execute(semver)
}
