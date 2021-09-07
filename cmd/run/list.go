package run

import (
	"context"
	"fmt"
	"sort"

	"github.com/logandavies181/tfd/cmd/config"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listRunCmd = &cobra.Command{
	Use:          "list",
	Aliases:      []string{"l"},
	Short:        "List runs",
	SilenceUsage: true,
	RunE:         listRun,
}

func init() {
	RunCmd.AddCommand(listRunCmd)

	listRunCmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to interact with")
}

type listRunConfig struct {
	*config.GlobalConfig

	Workspace string
}

func getListRunConfig(cmd *cobra.Command) (*listRunConfig, error) {
	viper.BindPFlags(cmd.Flags())

	gCfg, err := config.GetGlobalConfig()
	if err != nil {
		return nil, err
	}

	var lCfg listRunConfig
	err = viper.Unmarshal(&lCfg)
	if err != nil {
		return nil, err
	}

	lCfg.GlobalConfig = gCfg

	return &lCfg, nil
}

func listRun(cmd *cobra.Command, _ []string) error {
	cfg, err := getListRunConfig(cmd)
	if err != nil {
		return err
	}

	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	runList, err := cfg.Client.Runs.List(
		cfg.Ctx,
		workspace.ID,
		tfe.RunListOptions{})
	if err != nil {
		return err
	}

	// TODO: truncate the list and handle pagination
	runs := runList.Items
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
	runList, err := client.Runs.List(ctx, workspaceId, tfe.RunListOptions{})
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
