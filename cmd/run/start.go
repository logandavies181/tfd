package run

import (
	"github.com/logandavies181/tfd/v2/cmd/config"
	"github.com/logandavies181/tfd/v2/cmd/flags"

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

			AutoApply:            viper.GetBool("auto-apply"),
			ConfigurationVersion: viper.GetString("configuration-version"),
			FireAndForget:        viper.GetBool("fire-and-forget"),
			Message:              viper.GetString("message"),
			Refresh:              viper.GetBool("refresh"),
			RefreshOnly:          viper.GetBool("refresh-only"),
			Replace:              viper.GetStringSlice("replace"),
			Targets:              viper.GetStringSlice("targets"),
			Vars:                 viper.GetStringMapString("var"),
			Watch:                viper.GetBool("watch"),
			Workspace:            viper.GetString("workspace"),
		}

		return config.startRun(CREATE)
	},
}

func init() {
	RunCmd.AddCommand(runStartCmd)

	flags.AddAutoApplyFlag(runStartCmd)
	flags.AddConfigurationVersionFlag(runStartCmd)
	flags.AddFireAndForgetFlag(runStartCmd)
	flags.AddMessageFlag(runStartCmd)
	flags.AddRefreshFlag(runStartCmd)
	flags.AddRefreshOnlyFlag(runStartCmd)
	flags.AddReplaceFlag(runStartCmd)
	flags.AddTargetFlag(runStartCmd)
	flags.AddVarFlag(runStartCmd)
	flags.AddWatchFlag(runStartCmd)
	flags.AddWorkspaceFlag(runStartCmd)
}

type runStartConfig struct {
	config.Config

	AutoApply            bool   `mapstructure:"auto-apply"`
	ConfigurationVersion string `mapstructure:"configuration-version"`
	FireAndForget        bool   `mapstructure:"fire-and-forget"`
	IsDestroy            bool
	Message              string
	Refresh              bool
	RefreshOnly          bool `mapstructure:"refresh-only"`
	Replace              []string
	Targets              []string
	Vars                 map[string]string
	Watch                bool
	Workspace            string
}
