package run

import (
	"fmt"

	"github.com/logandavies181/tfd/v2/cmd/config"
	"github.com/logandavies181/tfd/v2/cmd/flags"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var applyRunCmd = &cobra.Command{
	Use:          "apply",
	Aliases:      []string{"a", "approve"},
	Short:        "Apply a run",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := applyRunConfig{
			Config: baseConfig,

			Watch:     viper.GetBool("watch"),
			Workspace: viper.GetString("workspace"),
		}

		return applyRun(config)
	},
}

func init() {
	RunCmd.AddCommand(applyRunCmd)

	flags.AddWatchFlag(applyRunCmd)
	flags.AddWorkspaceFlag(applyRunCmd)

	viper.BindPFlags(applyRunCmd.Flags())
}

type applyRunConfig struct {
	config.Config

	Watch     bool
	Workspace string
}

func applyRun(cfg applyRunConfig) error {
	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	if cfg.Watch {
		err := waitForWorkspaceToHaveCurrentRun(cfg.Ctx, cfg.Client, cfg.Org, cfg.Workspace)
		if err != nil {
			return err
		}

		run := workspace.CurrentRun
		return watchAndAutoApplyRun(cfg.Ctx, cfg.Client, cfg.Org, cfg.Workspace, run, true, cfg.Address)
	}

	runList, err := cfg.Client.Runs.List(
		cfg.Ctx,
		workspace.ID,
		&tfe.RunListOptions{})
	if err != nil {
		return err
	}

	var confirmableRunsList []string
	for _, r := range runList.Items {
		if r.Actions.IsConfirmable {
			confirmableRunsList = append(confirmableRunsList, r.ID)
		}
	}

	if len(confirmableRunsList) == 0 {
		return fmt.Errorf("No confirmable runs on workspace %s", cfg.Workspace)
	} else if len(confirmableRunsList) > 1 {
		return fmt.Errorf(
			"%d confirmable runs on workspace %s. Unsure how to proceed",
			len(confirmableRunsList),
			cfg.Workspace)
	}

	// confirm the run
	err = cfg.Client.Runs.Apply(cfg.Ctx, confirmableRunsList[0], tfe.RunApplyOptions{})
	if err != nil {
		return err
	}

	fmt.Println(confirmableRunsList[0])

	return nil
}
