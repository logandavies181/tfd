package run

import (
	"github.com/logandavies181/tfd/cmd/flags"

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

			AutoApply: viper.GetBool("auto-apply"),
			Message:   viper.GetString("message"),
			Watch:     viper.GetBool("watch"),
			Workspace: viper.GetString("workspace"),
		}

		return config.startRun(destroy)
	},
}

func init() {
	RunCmd.AddCommand(destroyRunCmd)

	flags.AddAutoApplyFlag(destroyRunCmd)
	flags.AddMessageFlag(destroyRunCmd)
	flags.AddWatchFlag(destroyRunCmd)
	flags.AddWorkspaceFlag(destroyRunCmd)

	viper.BindPFlags(destroyRunCmd.Flags())
}
