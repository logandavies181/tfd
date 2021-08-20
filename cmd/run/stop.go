package run

import (
	"context"
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/workspace"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stopRunCmd = &cobra.Command{
	Use:          "stop",
	Aliases:      []string{"s"},
	Short:        "Stop runs",
	SilenceUsage: true,
	RunE:         stopRun,
}

func init() {
	RunCmd.AddCommand(stopRunCmd)

	stopRunCmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to interact with")
}

type stopRunConfig struct {
	*config.GlobalConfig

	Watch     bool
	Workspace string
}

func getStopRunConfig(cmd *cobra.Command) (*stopRunConfig, error) {
	viper.BindPFlags(cmd.Flags())

	gCfg, err := config.GetGlobalConfig()
	if err != nil {
		return nil, err
	}

	var lCfg stopRunConfig
	err = viper.Unmarshal(&lCfg)
	if err != nil {
		return nil, err
	}

	lCfg.GlobalConfig = gCfg

	return &lCfg, nil
}

func stopRun(cmd *cobra.Command, _ []string) error {
	cfg, err := getStopRunConfig(cmd)
	if err != nil {
		return err
	}

	workspace, err := workspace.GetWorkspaceByName(*cfg.Client, cfg.Ctx, cfg.Org, cfg.Workspace)
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

	waitCounter := 0
	errChan := make(chan (error), len(runList.Items))
	var stopList []string

	// Cancel runs first, then discard runs
	for _, r := range runList.Items {
		if r.Actions.IsCancelable {
			stopList = append(stopList, r.ID)
			waitCounter++
			go discardOrCancelRun(cfg.Client, cfg.Ctx, r.ID, errChan, "Cancel")
		}
	}
	for waitCounter != 0 {
		err := <-errChan
		if err != nil {
			return err
		}
		waitCounter--
	}

	// Discard runs now that we've cancelled the rest
	for _, r := range runList.Items {
		if r.Actions.IsDiscardable {
			stopList = append(stopList, r.ID)
			waitCounter++
			go discardOrCancelRun(cfg.Client, cfg.Ctx, r.ID, errChan, "Discard")
		}
	}
	for waitCounter != 0 {
		err := <-errChan
		if err != nil {
			return err
		}
		waitCounter--
	}

	fmt.Printf("Stopped runs: %s\n", stopList)

	return nil
}

func discardOrCancelRun(client *tfe.Client, ctx context.Context, runId string, errChan chan (error), action string) {
	switch action {
	case "Discard":
		err := client.Runs.Discard(ctx, runId, tfe.RunDiscardOptions{})
		errChan <- err
	case "Cancel":
		err := client.Runs.Cancel(ctx, runId, tfe.RunCancelOptions{})
		errChan <- err
	default:
		errChan <- fmt.Errorf("Can only discard or cancel runs using discardOrCancelRun")
	}
}
