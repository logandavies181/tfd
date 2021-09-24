package plan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-tfe"
)

type PlanError struct {
	*tfe.Plan

	Message string
}

func (p PlanError) Error() string {
	return p.Message
}

func FormatPlanUrl(address string, plan *tfe.Plan) string {
	return fmt.Sprintf("%s/%s", address, plan.ID)
}

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
				return PlanError{Plan: p, Message: "Plan Errored"}
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
