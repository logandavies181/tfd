package vars

import (
	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"

	"github.com/spf13/cobra"
)

var varsSetCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := varsSetConfig{
			Config: baseConfig,
		}

		return varsSet(config)
	},
}

func init() {
	VarsCmd.AddCommand(varsSetCmd)
}

type varsSetConfig struct {
	config.Config
}

func varsSet(cfg varsSetConfig) error {
	return nil
}
