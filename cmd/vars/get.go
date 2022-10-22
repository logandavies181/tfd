package vars

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var varsGetCmd = &cobra.Command{
	Use:          "get",
	Aliases:      []string{"g"},
	Short:        "Get the value of a variable on a workspace",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := varsGetConfig{
			Config: baseConfig,

			Key:       viper.GetString("key"),
			Verbose:   viper.GetBool("verbose"),
			Workspace: viper.GetString("workspace"),
		}

		return varsGet(config)
	},
}

func init() {
	VarsCmd.AddCommand(varsGetCmd)

	flags.AddKeyFlag(varsGetCmd)
	flags.AddVerboseFlag(varsGetCmd)
	flags.AddWorkspaceFlag(varsGetCmd)
}

type varsGetConfig struct {
	config.Config

	Key       string
	Verbose   bool
	Workspace string
}

func varsGet(cfg varsGetConfig) error {
	wsVars, err := getAllVarsByWorkspaceName(cfg.Config, cfg.Workspace)
	if err != nil {
		return err
	}

	// The API sucks here so we need to search through the vars on our own
	var foundVar *tfe.Variable
	for _, v := range wsVars {
		if v.Key == cfg.Key {
			foundVar = v
		}
	}

	if foundVar == nil {
		return fmt.Errorf("Could not find var %s on workspace %s", cfg.Key, cfg.Workspace)
	}

	fmt.Println(formatVar(foundVar, cfg.Verbose))

	return nil
}
