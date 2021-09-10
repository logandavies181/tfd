package plan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-tfe"
)

func FormatResourceChanges(p *tfe.Plan) string {
	return fmt.Sprintf(
		"Plan: %d to add, %d to change, %d to destroy.",
		p.ResourceAdditions,
		p.ResourceChanges,
		p.ResourceDestructions)
}

// watchRun periodically checks the Run and returns when it is a finished, errored, or waiting for confirmation
func WatchPlan(ctx context.Context, client *tfe.Client, planId string) error {
	for {
		p, err := client.Plans.Read(ctx, planId)
		if err != nil {
			return err
		}

		if IsPlanFinished(p) {
			if p.Status == tfe.PlanErrored {
				return fmt.Errorf("Plan errored")
			}

			return nil
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

func IsPlanFinished(p *tfe.Plan) bool {
	switch p.Status {
	case tfe.PlanCanceled,
		tfe.PlanErrored,
		tfe.PlanFinished:

		return true
	default:
		return false
	}
}
