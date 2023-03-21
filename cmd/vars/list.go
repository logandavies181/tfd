package vars

import (
	"github.com/logandavies181/tfd/v2/cmd/config"
	"github.com/logandavies181/tfd/v2/cmd/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var varsListCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"l"},
	Short:        "List the variables for a workspace",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := varsListConfig{
			Config: baseConfig,

			Verbose:   viper.GetBool("verbose"),
			Workspace: viper.GetString("workspace"),
		}

		return varsList(config)
	},
}

func init() {
	VarsCmd.AddCommand(varsListCmd)

	flags.AddVerboseFlag(varsListCmd)
	flags.AddWorkspaceFlag(varsListCmd)
}

type varsListConfig struct {
	config.Config

	Verbose   bool
	Workspace string
}

func varsList(cfg varsListConfig) error {
	wsVars, err := getAllVarsByWorkspaceName(cfg.Config, cfg.Workspace)
	if err != nil {
		return err
	}

	printVarsTable(wsVars, cfg.Verbose)

	return nil
}
