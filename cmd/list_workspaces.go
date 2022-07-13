package cmd

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	"github.com/logandavies181/tfd/cmd/workspace"
	"github.com/logandavies181/tfd/pkg/pagination"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var listWorkspacesCmd = &cobra.Command{
	Use:          "list-workspaces",
	Aliases:      []string{"lw"},
	Short:        "List Terraform Cloud workspaces you have access to",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		return listWorkspaces(baseConfig)
	},
}

func init() {
	rootCmd.AddCommand(listWorkspacesCmd)
}

func listWorkspaces(cfg config.Config) error {
	var workspaces []*tfe.Workspace
	err := pagination.WithPagination(func(pg *tfe.Pagination) (bool, error) {
		workspaceListResp, err := cfg.Client.Workspaces.List(cfg.Ctx, cfg.Org, &tfe.WorkspaceListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: pg.NextPage,
			},
		})
		if err != nil {
			return false, err
		}
		if workspaceListResp.Pagination != nil {
			*pg = *workspaceListResp.Pagination
		}

		workspaces = append(workspaces, workspaceListResp.Items...)

		return false, nil
	})
	if err != nil {
		return err
	}

	workspace.SortWorkspacesByName(workspaces)
	for _, ws := range workspaces {
		fmt.Println(ws.Name)
	}

	return nil
}
