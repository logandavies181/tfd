package run

import (
	"context"
	"fmt"
	"time"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/plan"
	"github.com/logandavies181/tfd/cmd/workspace"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var runStartCmd = &cobra.Command{
	Use:          "start",
	Aliases:      []string{"s"},
	Short:        "Start a run",
	SilenceUsage: true,
	RunE:         runStart,
}

func init() {
	RunCmd.AddCommand(runStartCmd)

	runStartCmd.Flags().BoolP("auto-apply", "a", false, "Automatically apply the plan once finished")
	runStartCmd.Flags().BoolP("watch", "", false, "Wait for the run to finish")
	runStartCmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to interact with")
}

type runStartConfig struct {
	*config.GlobalConfig

	AutoApply bool `mapstructure:"auto-apply"`
	Watch     bool
	Workspace string
}

func getRunStartConfig(cmd *cobra.Command) (*runStartConfig, error) {
	viper.BindPFlags(cmd.Flags())

	gCfg, err := config.GetGlobalConfig()
	if err != nil {
		return nil, err
	}

	var lCfg runStartConfig
	err = viper.Unmarshal(&lCfg)
	if err != nil {
		return nil, err
	}

	lCfg.GlobalConfig = gCfg

	return &lCfg, nil
}

func runStart(cmd *cobra.Command, _ []string) error {
	cfg, err := getRunStartConfig(cmd)
	if err != nil {
		return err
	}

	workspace, err := workspace.GetWorkspaceByName(*cfg.Client, cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	r, err := cfg.Client.Runs.Create(
		cfg.Ctx,
		tfe.RunCreateOptions{
			Workspace: workspace,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println(r.ID)

	if cfg.Watch || cfg.AutoApply {
		err = watchAndAutoApplyRun(cfg.Ctx, cfg.Client, cfg.Org, workspace.Name, r, cfg.AutoApply)
		if err != nil {
			return err
		}
	}

	return nil
}

// watchAndAutoApplyRun waits for a run to plan and optionally auto-applies it, waiting for the apply to finish if so.
// It will return an error if it detects a queue on the workspace
func watchAndAutoApplyRun(ctx context.Context, client *tfe.Client, org, workspaceName string, r *tfe.Run, autoApply bool) error {
	// check if there's a queue
	err := waitForQueueStatus(ctx, client, org, workspaceName, r.ID)
	if err != nil {
		return err
	}

	fmt.Printf("Plan %s running. Waiting for it to finish..\n", r.Plan.ID)
	err = plan.WatchPlan(ctx, client, r.Plan.ID)
	if err != nil {
		return err
	}

	fmt.Println(plan.FormatResourceChanges(r.Plan))

	if autoApply {

		time.Sleep(1 * time.Second)

		for {
			r, err := client.Runs.Read(ctx, r.ID)
			if err != nil {
				return err
			}

			if isRunFinished(r) {
				fmt.Printf("Run %s finished with status: %s\n", r.ID, r.Status)
				return nil
			} else if r.Actions.IsConfirmable {
				err = client.Runs.Apply(ctx, r.ID, tfe.RunApplyOptions{})
				if err != nil {
					return err
				}
			} else {
				break
			}
			time.Sleep(5 * time.Second)
		}

		fmt.Println("Run confirmed")

		fmt.Println("Waiting for apply..")
		watchRun(ctx, client, r.ID)
		if err != nil {
			return err
		}

		r, err := client.Runs.Read(ctx, r.ID)
		if err != nil {
			return err
		}

		if isRunFinished(r) {
			fmt.Println("Run finished")

			fmt.Println(formatResourceChanges(r.Apply))
		}
	}

	return nil
}

// watchRun periodically checks the Run and returns when it is a finished, errored, or waiting for confirmation
func watchRun(ctx context.Context, client *tfe.Client, runId string) error {
	for {
		r, err := client.Runs.Read(ctx, runId)
		if err != nil {
			return err
		}

		if isRunFinished(r) {
			return nil
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

func isRunFinished(r *tfe.Run) bool {
	switch r.Status {
	case tfe.RunApplied,
		tfe.RunCanceled,
		tfe.RunDiscarded,
		tfe.RunErrored,
		tfe.RunPlannedAndFinished:

		return true
	default:
		return false
	}
}

// waitForQueueStatus periodically checks workspace.CurrentRun and returns once the current run is active. Err will be
// nil if the current run is the active one and non-nil if it is some other run
func waitForQueueStatus(ctx context.Context, client *tfe.Client, org, workspaceName, runId string) error {
	for {
		workspace, err := client.Workspaces.Read(ctx, org, workspaceName)
		if err != nil {
			return err
		}

		r, err := client.Runs.Read(ctx, runId)
		if err != nil {
			return err
		}

		if workspace.CurrentRun == nil {
			time.Sleep(5 * time.Second)
		} else if workspace.CurrentRun.ID != r.ID {
			if !isRunFinished(workspace.CurrentRun) {
				// Current run is someone else. Don't wait for queue, just exit
				return fmt.Errorf("Workspace is currently locked by %s. "+
					"Complete or discard that run before attempting to queue",
					workspace.CurrentRun.ID)
			} else {
				// Current run isn't running and isn't us. Wait for Terraform Cloud to catch up
				time.Sleep(5 * time.Second)
			}
		} else {
			// We're the current run. Return now
			return nil
		}
	}
}
