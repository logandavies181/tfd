package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/logandavies181/tfd/cmd"
)

func main() {
	debugInfo, ok := debug.ReadBuildInfo()
	if !ok {
		fmt.Fprintln(os.Stderr, "Could not determine build info. Your tfd version is corrupt")
		os.Exit(1)
	}

	version := debugInfo.Main.Version
	if version == "(devel)" {
		version = "v0.0.0+devel"
	}

	cmd.Execute(version)
}
