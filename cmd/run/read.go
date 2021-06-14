package run

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"

	"gopkg.in/yaml.v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var readRuncmd = &cobra.Command{
	Use:          "read",
	Aliases:      []string{"r"},
	Short:        "Read a run",
	SilenceUsage: true,
	RunE:         readRun,
}

func init() {
	RunCmd.AddCommand(readRuncmd)

	readRuncmd.Flags().StringP("run-id", "r", "", "Run ID to read")
	readRuncmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to read the run from")
}

type readRunConfig struct {
	*config.GlobalConfig

	RunId string `mapstructure:"run-id"`
	Workspace string
}

func getReadRunConfig(cmd *cobra.Command) (*readRunConfig, error) {
	viper.BindPFlags(cmd.Flags())

	gCfg, err := config.GetGlobalConfig()
	if err != nil {
		return nil, err
	}

	var lCfg readRunConfig
	err = viper.Unmarshal(&lCfg)
	if err != nil {
		return nil, err
	}

	lCfg.GlobalConfig = gCfg

	return &lCfg, nil
}

func readRun(cmd *cobra.Command, _ []string) error {
	cfg, err := getReadRunConfig(cmd)
	if err != nil {
		return err
	}

	var runId string
	if cfg.RunId == "" {
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
