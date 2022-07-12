package cv

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	//"github.com/logandavies181/tfd/pkg/pagination"

	//"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cvListCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"l"},
	Short:        "List configuration versions",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := cvListConfig{
			Config: baseConfig,

			Workspace: viper.GetString("workspace"),
		}

		return cvList(config)
	},
}

func init() {
	CvCmd.AddCommand(cvListCmd)

	flags.AddWorkspaceFlag(cvListCmd)
}

type cvListConfig struct {
	config.Config

	Workspace string
}

func cvList(cfg cvListConfig) error {
	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return fmt.Errorf("Could not read workspace Id: %v", err)
	}

	cvl, err := cfg.Client.ConfigurationVersions.List(cfg.Ctx, workspace.ID, nil)
	if err != nil {
		return fmt.Errorf("Could not list configuration versions: %v", err)
	}

	for _, v := range cvl.Items {
		fmt.Println(v.ID)
	}

	return nil
}
