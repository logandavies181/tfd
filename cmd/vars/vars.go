package vars

import (
	"github.com/spf13/cobra"
)

var VarsCmd = &cobra.Command{
	Use:           "variables",
	Aliases:       []string{"vars", "v"},
	Short:         "Commands for interacting with Workspace Variables",
	SilenceErrors: true,
}
