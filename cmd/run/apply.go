package run

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/workspace"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var applyRunCmd = &cobra.Command{
	Use:          "apply",
	Aliases:      []string{"a", "approve"},
	Short:        "Apply a run",
	SilenceUsage: true,
	RunE:         applyRun,
}

func init() {
	RunCmd.AddCommand(applyRunCmd)

	applyRunCmd.Flags().BoolP("watch", "", false, "Wait for the run to finish")
	applyRunCmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to interact with")
}

type applyRunConfig struct {
	*config.GlobalConfig

	Watch     bool
	Workspace string
}

func getApplyRunConfig(cmd *cobra.Command) (*applyRunConfig, error) {
	viper.BindPFlags(cmd.Flags())

	gCfg, err := config.GetGlobalConfig()
	if err != nil {
		return nil, err
	}

	var lCfg applyRunConfig
	err = viper.Unmarshal(&lCfg)
	if err != nil {
		return nil, err
	}

	lCfg.GlobalConfig = gCfg

	return &lCfg, nil
}

func applyRun(cmd *cobra.Command, _ []string) error {
	cfg, err := getApplyRunConfig(cmd)
	if err != nil {
		return err
	}

	workspace, err := workspace.GetWorkspaceByName(*cfg.Client, cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	runList, err := cfg.Client.Runs.List(
		cfg.Ctx,
		workspace.ID,
		tfe.RunListOptions{})
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
