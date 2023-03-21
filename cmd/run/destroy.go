package run

import (
	"github.com/logandavies181/tfd/v2/cmd/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var destroyRunCmd = &cobra.Command{
	Use:          "destroy",
	Aliases:      []string{"d"},
	Short:        "Start a destroy run",
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

		return config.startRun(DESTROY)
	},
}

func init() {
	RunCmd.AddCommand(destroyRunCmd)

	flags.AddAutoApplyFlag(destroyRunCmd)
	flags.AddConfigurationVersionFlag(destroyRunCmd)
	flags.AddFireAndForgetFlag(destroyRunCmd)
	flags.AddMessageFlag(destroyRunCmd)
	flags.AddRefreshFlag(destroyRunCmd)
	flags.AddRefreshOnlyFlag(destroyRunCmd)
	flags.AddReplaceFlag(destroyRunCmd)
	flags.AddTargetFlag(destroyRunCmd)
	flags.AddVarFlag(destroyRunCmd)
	flags.AddWatchFlag(destroyRunCmd)
	flags.AddWorkspaceFlag(destroyRunCmd)
}
