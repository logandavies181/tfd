package flags

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("path", "p", "", "Path to Terraform Directory")
}

func AddWorkspaceFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to upload to")

	flagValidations = append(flagValidations, func() error {
		if viper.GetString("workspace") == "" {
			return fmt.Errorf("workspace must be set")
		}
		return nil
	})
}

func AddNoUpdateWorkingdirFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("no-update-workdingir", "d", false,
		"Skip updating the Terraform Working Directory for the workspace")
}

func AddAutoApplyFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("auto-apply", "a", false, "Automatically apply the plan once finished")
}

func AddWatchFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("watch", "", false, "Wait for the run/plan to finish")
}

func AddRunIdFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("run-id", "r", "", "Run ID to read")
}
