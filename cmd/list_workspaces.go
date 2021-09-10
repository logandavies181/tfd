package cmd

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	"github.com/logandavies181/tfd/cmd/workspace"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var listWorkspacesCmd = &cobra.Command{
	Use:          "list-workspaces",
	Aliases:      []string{"lw"},
	Short:        "List Terraform Cloud workspaces you have access to",
	SilenceUsage: true,
	RunE: func(_ *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd()
		if err != nil {
			return err
		}

		return listWorkspaces(baseConfig)
	},
}

func init() {
	rootCmd.AddCommand(listWorkspacesCmd)
}

func listWorkspaces(cfg *config.Config) error {
	var workspaces []*tfe.Workspace
	pagination := &tfe.Pagination{
		NextPage:   1,
		TotalPages: -1,
	}
	for {
		if pagination == nil || pagination.CurrentPage == pagination.TotalPages {
			break
		}
		workspaceListResp, err := cfg.Client.Workspaces.List(cfg.Ctx, cfg.Org, tfe.WorkspaceListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: pagination.NextPage,
			},
		})
		if err != nil {
			return err
		}
		pagination = workspaceListResp.Pagination

		workspaces = append(workspaces, workspaceListResp.Items...)
	}

	workspace.SortWorkspacesByName(workspaces)
	for _, ws := range workspaces {
		fmt.Println(ws.Name)
	}

	return nil
}
