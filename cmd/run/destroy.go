package run

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var destroyRunCmd = &cobra.Command{
	Use:          "destroy",
	Aliases:      []string{"d"},
	Short:        "Start a destroy run",
	SilenceUsage: true,
	RunE:         destroyRun,
}

func init() {
	RunCmd.AddCommand(destroyRunCmd)

	destroyRunCmd.Flags().BoolP("auto-apply", "a", false, "Automatically apply the plan once finished")
	destroyRunCmd.Flags().BoolP("watch", "", false, "Wait for the run to finish")
	destroyRunCmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to interact with")
}

type destroyRunConfig struct {
	*config.GlobalConfig

	AutoApply bool `mapstructure:"auto-apply"`
	Watch     bool
	Workspace string
}

func getDestroyRunConfig(cmd *cobra.Command) (*destroyRunConfig, error) {
	viper.BindPFlags(cmd.Flags())

	gCfg, err := config.GetGlobalConfig()
	if err != nil {
		return nil, err
	}

	var lCfg destroyRunConfig
	err = viper.Unmarshal(&lCfg)
	if err != nil {
		return nil, err
	}

	lCfg.GlobalConfig = gCfg

	return &lCfg, nil
}

func destroyRun(cmd *cobra.Command, _ []string) error {
	cfg, err := getDestroyRunConfig(cmd)
	if err != nil {
		return err
	}

	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	isDestroy := true

	r, err := cfg.Client.Runs.Create(
		cfg.Ctx,
		tfe.RunCreateOptions{
			IsDestroy: &isDestroy,
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
