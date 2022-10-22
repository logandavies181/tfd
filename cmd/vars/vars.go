package vars

import (
	"fmt"
	"os"
	"strings"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/pkg/pagination"

	"github.com/hashicorp/go-tfe"
	"github.com/olekukonko/tablewriter"
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

func printVarsTable(vars []*tfe.Variable, verbose bool) {
	table := tablewriter.NewWriter(os.Stdout)

	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)

	for _, v := range vars {
		varString := formatVar(v, verbose)

		// chuck v.Key at the front
		cols := append([]string{v.Key}, strings.Split(varString, "\t")...)

		table.Append(cols)
	}

	table.Render()
}

func formatVar(v *tfe.Variable, verbose bool) string {
	if v == nil {
		return ""
	}

	fields := []string{
		v.Value,
	}

	if verbose {
		fields = append(
			fields,
			formatDescription(v.Description),
			formatIsHCL(v.HCL),
			formatIsSensitive(v.Sensitive),
			string(v.Category),
		)
	}

	return strings.Join(fields, "\t")
}

func formatDescription(desc string) string {
	if desc == "" {
		return "no description"
	}

	return desc
}

func formatIsHCL(isHCL bool) string {
	if isHCL {
		return "hcl"
	}

	return "normal"
}

func formatIsSensitive(isSensitive bool) string {
	if isSensitive {
		return "sensitive"
	}

	return "not sensitive"
}
