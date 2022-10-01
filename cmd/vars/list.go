package vars

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	"github.com/logandavies181/tfd/pkg/pagination"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var varsListCmd = &cobra.Command{
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := varsListConfig{
			Config: baseConfig,

			Workspace: viper.GetString("workspace"),
		}

		return varsList(config)
	},
}

func init() {
	VarsCmd.AddCommand(varsListCmd)

	flags.AddWorkspaceFlag(varsListCmd)
}

type varsListConfig struct {
	config.Config

	Workspace string
}

func varsList(cfg varsListConfig) error {
	ws, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	var wsVars []*tfe.Variable
	err = pagination.WithPagination(func(pg *tfe.Pagination) (bool, error) {
		varsListResp, err := cfg.Client.Variables.List(cfg.Ctx, ws.ID, &tfe.VariableListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: pg.NextPage,
			},
		})
		if err != nil {
			return false, err
		}
		if varsListResp.Pagination != nil {
			*pg = *varsListResp.Pagination
		}

		wsVars = append(wsVars, varsListResp.Items...)

		return false, nil
	})
	if err != nil {
		return err
	}

	for _, v := range wsVars {
		fmt.Printf("%s: %s", v.Key, v.Value)
	}

	return nil
}
