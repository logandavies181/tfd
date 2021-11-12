package plan

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-tfe"
)

// overridden in tests to speed them up
var pollingInterval time.Duration = 10

type PlanError struct {
	*tfe.Plan

	Message string
}

func (p PlanError) Error() string {
	return p.Message
}

func FormatResourceChanges(p *tfe.Plan) string {
	return fmt.Sprintf(
		"Plan: %d to add, %d to change, %d to destroy.",
		p.ResourceAdditions,
		p.ResourceChanges,
		p.ResourceDestructions)
}

// WatchRun periodically checks the Run and returns when it is finished, errored, or waiting for confirmation
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
			time.Sleep(pollingInterval * time.Second)
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
