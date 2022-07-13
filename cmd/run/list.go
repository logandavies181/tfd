package run

import (
	"context"
	"fmt"
	"sort"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	"github.com/logandavies181/tfd/pkg/pagination"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listRunCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"l"},
	Short:        "List runs",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := &listRunConfig{
			Config: baseConfig,

			MaxItems: viper.GetInt("max-items"),
			Workspace: viper.GetString("workspace"),
		}

		return listRun(config)
	},
}

func init() {
	RunCmd.AddCommand(listRunCmd)

	flags.AddMaxItemsFlag(listRunCmd)
	flags.AddWorkspaceFlag(listRunCmd)
}

type listRunConfig struct {
	config.Config

	MaxItems  int `mapstructure:"max-items"`
	Workspace string
}

func listRun(cfg *listRunConfig) error {
	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	var runs []*tfe.Run
	err = pagination.WithPagination(func(pg *tfe.Pagination) (bool, error) {
		runList, err := cfg.Client.Runs.List(
			cfg.Ctx,
			workspace.ID,
			&tfe.RunListOptions{})
		if err != nil {
			return false, err
		}

		runs = append(runs, runList.Items...)
		if len(runs) >= cfg.MaxItems {
			runs = runs[:cfg.MaxItems]
			return true, nil
		}

		*pg = *runList.Pagination

		return false, nil
	})
	if err != nil {
		return err
	}

	sortRunsByCreateTime(runs)
	for _, r := range runs {
		fmt.Printf("%s\t%s\t%s\n", r.CreatedAt.Format("Jan 2 15:04:05"), r.ID, r.Status)
	}

	return nil
}

type RunTimeSorter []*tfe.Run

func (rts RunTimeSorter) Len() int {
	return len(rts)
}

func (rts RunTimeSorter) Less(i, j int) bool {
	return rts[i].CreatedAt.Before(rts[j].CreatedAt)
}

func (rts RunTimeSorter) Swap(i, j int) {
	rts[i], rts[j] = rts[j], rts[i]
}

func sortRunsByCreateTime(runs []*tfe.Run) {
	sort.Sort(RunTimeSorter(runs))
}

func getConfirmableRunByWorkspaceId(client *tfe.Client, ctx context.Context, workspaceId string) (string, error) {
	runList, err := client.Runs.List(ctx, workspaceId, &tfe.RunListOptions{})
	if err != nil {
		return "", err
	}

	for _, r := range runList.Items {
		if r.Actions.IsConfirmable {
			return r.ID, nil
		}
	}

	return "", fmt.Errorf("No confirmable Runs on workspace %s", workspaceId)
}
