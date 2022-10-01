package vars

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/pkg/pagination"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

var VarsCmd = &cobra.Command{
	Use:           "variables",
	Aliases:       []string{"vars", "v"},
	Short:         "Commands for interacting with Workspace Variables",
	SilenceErrors: true,
}

func categoryType(s string) *tfe.CategoryType {
	ct := tfe.CategoryType(s)

	return &ct
}

func getVarByName(cfg config.Config, workspaceName string, varName string) (*tfe.Variable, error) {
	wsVars, err := getAllVarsByWorkspaceName(cfg, workspaceName)
	if err != nil {
		return nil, err
	}

	for _, wsVar := range wsVars {
		if wsVar != nil && wsVar.Key == varName {
			return wsVar, nil
		}
	}

	return nil, fmt.Errorf("Could not find var with key: %s", varName)
}

func getAllVarsByWorkspaceName(cfg config.Config, workspaceName string) ([]*tfe.Variable, error) {
	ws, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, workspaceName)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return wsVars, nil
}
