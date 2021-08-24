package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/logandavies181/tfd/cmd"
)

var version string

func main() {
	cmd.Execute(buildVersion())
}

// buildVersion checks if version has been set at build time, otherwise uses debug.ReadBuildInfo to infer a tag.
// debug.ReadBuildInfo behaves differently given the following scenarios so buildVersion does the following:
// Locally compiled - Main.Version = "devel"; return "0.0.0+devel"
// go get tfd@<some git hash> - Main.Version = "0.0.0+v0.0.0-<timestamp>-<short_hash>"; return 0.0.0-<timestamp>-<short_hash>
// go get tfd@<tag> - Main.Version = "0.0.0-v<tag>"; return "<tag>"
func buildVersion() string {
	if version == "" {
		debugInfo, ok := debug.ReadBuildInfo()
		if !ok {
			fmt.Fprintln(os.Stderr, "Could not determine build info. Your tfd version is corrupt")
			os.Exit(1)
		}

		version := debugInfo.Main.Version
		if version == "(devel)" {
			return "0.0.0+devel"
		}

		foundVersionStart := false
		version = strings.TrimLeftFunc(version, func(r rune) bool {
			if r == 'v' {
				// only false for first v - trim it but not the others
				if foundVersionStart {
					return false
				}
				foundVersionStart = true

			}
			return !foundVersionStart
		})
	}
	return version
}
