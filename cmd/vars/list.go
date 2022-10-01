package vars

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"

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

			Workspace: viper.GetString("workspace"),
		}

		return varsList(config)
	},
}

func init() {
	VarsCmd.AddCommand(varsListCmd)

	flags.AddWorkspaceFlag(varsListCmd)
}

type varsListConfig struct {
	config.Config

	Workspace string
}

func varsList(cfg varsListConfig) error {
	wsVars, err := getAllVarsByWorkspaceName(cfg.Config, cfg.Workspace)
	if err != nil {
		return err
	}

	for _, v := range wsVars {
		fmt.Printf("%s: %s\n", v.Key, v.Value)
	}

	return nil
}
