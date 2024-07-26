package flags

import (
	"fmt"
	"github.com/pkg/errors"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func AddAutoApplyFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("auto-apply", "a", false, "Automatically apply the plan once finished")
}

func AddCategoryFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("category", "", "", `Variable category. Acceptable values are "env", "policy-set", and "terraform"`)

	addValidation(cmd.Name(), func() error {
		switch cat := viper.GetString("category"); cat {
		case "env", "policy-set", "terraform":
		case "":
			viper.Set("category", "terraform")
		default:
			return fmt.Errorf(`Category must be one of: "env", "policy-set", or "terraform"`)
		}
		return nil
	})
}

func AddConfigurationVersionFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("configuration-version", "", "", "Configuration version to create a run against")
}

func AddDescriptionFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("description", "", "", "Description")
}

func AddFireAndForgetFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("fire-and-forget", "f", false,
		"Non-interactively apply the plan once finished. Warning: this will still auto apply even if tfd exits. Use --auto-apply instead for safety")
}

func AddKeyFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("key", "", "", "Variable key name")

	addValidation(cmd.Name(), func() error {
		if viper.GetString("key") == "" {
			return fmt.Errorf("Key flag must be set")
		}

		return nil
	})
}

func AddHclFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("hcl", "", false, "Declare an HCL type variable")
}

func AddMaxItemsFlag(cmd *cobra.Command) {
	cmd.Flags().UintP("max-items", "", 10, "Max number of items to fetch")
}

func AddMessageFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("message", "m", "", "Specifies the reason for the current action")
}

func AddNoClobberFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("no-clobber", "", false, "Don't override existing")
}

func AddPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("path", "p", "", "Path to project. Can be any subdirectory of the workspace root")
}

func AddRootPathFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("rootpath", "", "", "Path to workspace root. Defaults to the folder detected by git as the workspace root")

	addValidation(cmd.Name(), func() error {
		rootPath := viper.GetString("rootpath")
		if rootPath == "" {
			return nil
		}
		stat, err := os.Stat(rootPath)
		if err != nil {
			return errors.Wrapf(err, "tfd: error accessing rootPath")
		}
		if !stat.IsDir() {
			return errors.Errorf("tfd: rootPath '%s' is not a directory", rootPath)
		}
		return nil
	})
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

func AddSensitiveFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("sensitive", "", false, "Whether or not this is a sensitive variable")
}

func AddTargetFlag(cmd *cobra.Command) {
	cmd.Flags().StringSliceP("target", "", []string{},
		"EXPERIMENTAL: Specifies the list of target addresses to use for the run. Not recommended for regular use. Can be specified many times")
}

func AddValueFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("value", "", "", "Variable value")
}

func AddVarFlag(cmd *cobra.Command) {
	cmd.Flags().StringToStringP("var", "", make(map[string]string),
		"EXPERIMENTAL: Sets variables for the current run, taking precedence over those set on the workspace. Can be specified many times")
}

func AddVerboseFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("verbose", "", false, "Print verbose output")
}

func AddWatchFlag(cmd *cobra.Command) {
	cmd.Flags().BoolP("watch", "", false, "Wait for the run/plan to finish")
}

func AddWorkspaceFlag(cmd *cobra.Command) {
	cmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to interact with")

	addValidation(cmd.Name(), func() error {
		if viper.GetString("workspace") == "" {
			return fmt.Errorf("workspace must be set")
		}
		return nil
	})
}
