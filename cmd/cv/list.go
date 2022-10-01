package cv

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	"github.com/logandavies181/tfd/pkg/pagination"

	"github.com/hashicorp/go-tfe"
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

			MaxItems:  viper.GetInt("max-items"),
			Workspace: viper.GetString("workspace"),
		}

		return cvList(config)
	},
}

func init() {
	CvCmd.AddCommand(cvListCmd)

	flags.AddMaxItemsFlag(cvListCmd)
	flags.AddWorkspaceFlag(cvListCmd)
}

type cvListConfig struct {
	config.Config

	MaxItems  int `mapstructure:"max-items"`
	Workspace string
}

func cvList(cfg cvListConfig) error {
	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return fmt.Errorf("Could not read workspace Id: %v", err)
	}

	var cvs []*tfe.ConfigurationVersion
	err = pagination.WithPagination(func(pg *tfe.Pagination) (bool, error) {
		cvl, err := cfg.Client.ConfigurationVersions.List(cfg.Ctx, workspace.ID, nil)
		if err != nil {
			return false, fmt.Errorf("Could not list configuration versions: %v", err)
		}

		cvs = append(cvs, cvl.Items...)

		if len(cvs) >= cfg.MaxItems {
			cvs = cvs[:cfg.MaxItems]
			return true, nil
		}

		*pg = *cvl.Pagination

		return false, nil
	})
	if err != nil {
		return err
	}

	for _, c := range cvs {
		fmt.Println(c.ID)
	}

	return nil
}
