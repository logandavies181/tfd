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

func AddMaxItemsFlag(cmd *cobra.Command) {
	cmd.Flags().UintP("max-items", "", 10, "Max number of items to fetch")
}

func AddMessageFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("message", "m", "", "Specifies the reason for the current action")
}

func AddPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("path", "p", "", "Path to project. Can be any subdirectory of the project, but it must be a git project")
}

func AddRefreshFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("refresh", "", true, "Determines if the run should check the state of created resources before creating plan")
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

func AddTargetFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("target", "", []string{},
		"EXPERIMENTAL: Specifies the list of target addresses to use for the run. Not recommended for regular use. Can be specified many times")
}

func AddVarFlag(cmd *cobra.Command) {
	cmd.Flags().StringToStringP("var", "", make(map[string]string),
		"EXPERIMENTAL: Sets variables for the current run, taking precedence over those set on the workspace. Can be specified many times")
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
