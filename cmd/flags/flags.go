package flags

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddAutoApplyFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("auto-apply", "a", false, "Automatically apply the plan once finished")
}

func AddConfigurationVersionFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("configuration-version", "", "", "Configuration version to create a run against")
}

func AddFireAndForgetFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("fire-and-forget", "f", false,
		"Non-interactively apply the plan once finished. Warning: this will still auto apply even if tfd exits. Use --auto-apply instead for safety")
}

func AddMessageFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("message", "m", "", "Specifies the reason for the current action")
}

func AddPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("path", "p", "", "Path to project. Can be any subdirectory of the project, but it must be a git project")
}

func AddRefreshFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("refresh", "", true, "Determines if the run should update the state prior to checking for differences")
}

func AddRefreshOnlyFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("refresh-only", "", false, "Determines whether the run should ignore config changes and refresh the state only")
}

func AddReplaceFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("replace", "", []string{},
		"EXPERIMENTAL: Specifies a list of addresses to recreate in the run. Not recommended for regular use. Can be specified many times")
}

func AddRunIdFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("run-id", "r", "", "Run ID to read")
}

func AddTargetsFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("targets", "", []string{},
		"EXPERIMENTAL: Specifies the list of target addresses to use for the run. Not recommended for regular use. Can be specified many times")
}

func AddWatchFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("watch", "", false, "Wait for the run/plan to finish")
}

func AddWorkspaceFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to upload to")

	addValidation(cmd.Name(), func() error {
		if viper.GetString("workspace") == "" {
			return fmt.Errorf("workspace must be set")
		}
		return nil
	})
}
