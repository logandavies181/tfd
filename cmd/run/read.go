package run

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

var readRuncmd = &cobra.Command{
	Use:          "read",
	Aliases:      []string{"r", "status"},
	Short:        "Read a run",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := &readRunConfig{
			Config: baseConfig,

			RunId:     viper.GetString("run-id"),
			Workspace: viper.GetString("workspace"),
		}

		return readRun(config)
	},
}

func init() {
	RunCmd.AddCommand(readRuncmd)

	flags.AddRunIdFlag(readRuncmd)
	flags.AddWorkspaceFlag(readRuncmd)
}

type readRunConfig struct {
	config.Config

	RunId     string
	Workspace string
}

func readRun(cfg *readRunConfig) error {
	var runId string
	if cfg.RunId == "" {
		var err error
		runId, err = getCurrentRun(cfg.Ctx, cfg.Client, cfg.Org, cfg.Workspace)
		if err != nil {
			return err
		}
	} else {
		runId = cfg.RunId
	}

	r, err := cfg.Client.Runs.Read(cfg.Ctx, runId)
	if err != nil {
		return err
	}

	runYaml, err := yaml.Marshal(r)
	if err != nil {
		return err
	}

	fmt.Println(string(runYaml))

	return nil
}
