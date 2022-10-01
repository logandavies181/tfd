package vars

import (
	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var VarsCmd = &cobra.Command{
	Use:           "variables",
	Aliases:       []string{"vars", "v"},
	Short:         "Commands for interacting with Workspace Variables",
	SilenceErrors: true,
}

func categoryType(s string) *tfe.CategoryType {
	ct := tfe.CategoryType(s)

	return &ct
}
