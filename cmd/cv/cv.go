package cv

import (
	"github.com/spf13/cobra"
)

var CvCmd = &cobra.Command{
	Use:           "configuration-version",
	Aliases:       []string{"cv"},
	Short:         "Commands for interacting with Configuration Versions",
	SilenceErrors: true,
}
