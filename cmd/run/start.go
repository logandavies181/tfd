package run

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runStartCmd = &cobra.Command{
	Use:          "start",
	Aliases:      []string{"s"},
	Short:        "Start a run",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := &runStartConfig{
			Config: baseConfig,

			AutoApply: viper.GetBool("auto-apply"),
			Watch:     viper.GetBool("watch"),
			Workspace: viper.GetString("workspace"),
		}

		return runStart(config)
	},
}

func init() {
	RunCmd.AddCommand(runStartCmd)

	flags.AddAutoApplyFlag(runStartCmd)
	flags.AddWatchFlag(runStartCmd)
	flags.AddWorkspaceFlag(runStartCmd)
}

type runStartConfig struct {
	*config.Config

	AutoApply bool `mapstructure:"auto-apply"`
	Watch     bool
	Workspace string
}

func runStart(cfg *runStartConfig) error {
	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	r, err := cfg.Client.Runs.Create(
		cfg.Ctx,
		tfe.RunCreateOptions{
			Workspace: workspace,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(r.ID)

	if cfg.Watch || cfg.AutoApply {
		err = watchAndAutoApplyRun(cfg.Ctx, cfg.Client, cfg.Org, workspace.Name, r, cfg.AutoApply)
		if err != nil {
			return err
		}
	}

	return nil
}
