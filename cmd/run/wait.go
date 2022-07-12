package run

import (
	"context"
	"fmt"
	"time"

	"github.com/logandavies181/tfd/cmd/plan"

	"github.com/hashicorp/go-tfe"
)

// watchAndAutoApplyRun waits for a run to plan and optionally auto-applies it, waiting for the apply to finish if so.
// It will return an error if it detects a queue on the workspace
func watchAndAutoApplyRun(ctx context.Context, client *tfe.Client, org, workspaceName string, r *tfe.Run, autoApply bool, address string) error {
	if r == nil {
		return fmt.Errorf("Fatal: Run is nil")
	}

	runUrl, err := FormatRunUrl(address, org, workspaceName, r.ID)
	if err != nil {
		return err
	}
	fmt.Println("View the plan in the UI:", runUrl)

	// check if there's a queue
	err = waitForQueueStatus(ctx, client, org, workspaceName, r.ID)
	if err != nil {
		return err
	}

	// r.Plan seems to be nil when we get it from the current workspace??
	var planId string
	if r.Plan == nil {
		run, err := client.Runs.Read(ctx, r.ID)
		if err != nil {
			return err
		}

		planId = run.Plan.ID
	} else {
		planId = r.Plan.ID
	}
	fmt.Printf("Plan %s running. Waiting for it to finish..\n", planId)

	err = plan.WatchPlan(ctx, client, planId)
	if err != nil {
		return err
	}

	// read the plan directly as the relation might be nil
	p, err := client.Plans.Read(ctx, planId)
	if err != nil {
		return err
	}
	fmt.Println(plan.FormatResourceChanges(p))

	if autoApply {

		time.Sleep(1 * time.Second)

		// Wait for run to be confirmable. TODO: check if this works with auto-apply enabled on the workspace
		for {
			r, err := client.Runs.Read(ctx, r.ID)
			if err != nil {
				return err
			}

			if isRunFinished(r) {
				if r.Status == tfe.RunErrored {
					return fmt.Errorf("Run errored")
				}
				fmt.Printf("Run %s finished with status: %s\n", r.ID, r.Status)
				return nil
			} else if r.Actions.IsConfirmable {
				err = client.Runs.Apply(ctx, r.ID, tfe.RunApplyOptions{})
				if err != nil {
					return err
				}
			} else if isRunWaitingBetweenPlanAndApplying(r) {
				// spin again
				continue
			} else {
				break
			}
			time.Sleep(5 * time.Second)
		}

		fmt.Println("Run confirmed")

		fmt.Println("Waiting for apply..")
		err := watchRun(ctx, client, r.ID)
		if err != nil {
			return err
		}

		r, err := client.Runs.Read(ctx, r.ID)
		if err != nil {
			return err
		}

		if isRunFinished(r) {
			fmt.Println("Run finished")

			// Try getting Apply directly instead of using relation, which may be nil
			appl, err := client.Applies.Read(ctx, r.Apply.ID)
			if err != nil {
				return err
			}
			fmt.Println(formatResourceChanges(appl))
		}
	}

	return nil
}

// watchRun periodically checks the Run and returns when it is finished, errored, or waiting for confirmation
func watchRun(ctx context.Context, client *tfe.Client, runId string) error {
	for {
		r, err := client.Runs.Read(ctx, runId)
		if err != nil {
			return err
		}

		if isRunFinished(r) {
			if r.Status == tfe.RunErrored {
				return fmt.Errorf("Run errored")
			}
			return nil
		} else {
			time.Sleep(5 * time.Second)
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

func isRunWaitingBetweenPlanAndApplying(r *tfe.Run) bool {
	switch r.Status {
	case tfe.RunConfirmed,
		tfe.RunCostEstimated,
		tfe.RunCostEstimating,
		tfe.RunPending,
		tfe.RunPolicyChecked,
		tfe.RunPolicyChecking,
		tfe.RunPolicyOverride,
		tfe.RunPolicySoftFailed:

		return true
	default:
		return false
	}
}

// waitForQueueStatus periodically checks workspace.CurrentRun and returns once the current run is active. Err will be
// nil if the current run is the active one and non-nil if it is some other run
func waitForQueueStatus(ctx context.Context, client *tfe.Client, org, workspaceName, runId string) error {
	for {
		err := waitForWorkspaceToHaveCurrentRun(ctx, client, org, workspaceName)
		if err != nil {
			return err
		}

		workspace, err := client.Workspaces.Read(ctx, org, workspaceName)
		if err != nil {
			return err
		}

		r, err := client.Runs.Read(ctx, runId)
		if err != nil {
			return err
		}

		if workspace.CurrentRun.ID != r.ID {
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

func waitForWorkspaceToHaveCurrentRun(ctx context.Context, client *tfe.Client, org, workspaceName string) error {
	for {
		workspace, err := client.Workspaces.Read(ctx, org, workspaceName)
		if err != nil {
			return err
		}

		if workspace.CurrentRun == nil {
			time.Sleep(5 * time.Second)
		} else {
			return nil
		}
	}
}
