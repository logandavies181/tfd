package plan

import (
	"testing"

	"github.com/logandavies181/tfd/mocks"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestFormatResourceChanges(t *testing.T) {
	plan := &tfe.Plan{
		ResourceAdditions:    1,
		ResourceChanges:      2,
		ResourceDestructions: 3,
	}

	resourceChangesStr := FormatResourceChanges(plan)

	assert.Equal(
		t,
		"Plan: 1 to add, 2 to change, 3 to destroy.",
		resourceChangesStr)
}

func TestIsPlanFinished(t *testing.T) {
	// expect true
	plan := &tfe.Plan{Status: tfe.PlanCanceled}
	isPlanFinished := IsPlanFinished(plan)
	assert.True(t, isPlanFinished)

	// expect false
	plan = &tfe.Plan{Status: "some other status"}
	isPlanFinished = IsPlanFinished(plan)
	assert.False(t, isPlanFinished)
}

func TestWatchPlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockPlans(ctrl)

	gomock.InOrder(
		m.EXPECT().
			Read(
				gomock.Any(),
				gomock.Eq("plan-123"),
			).
			Return(
				&tfe.Plan{
					Status: "some status",
				},
				nil,
			),
		m.EXPECT().
			Read(
				gomock.Any(),
				gomock.Eq("plan-123"),
			).
			Return(
				&tfe.Plan{
					Status: tfe.PlanFinished,
				},
				nil,
			),
	)

	cfg := mocks.MockConfig()
	cfg.Client.Plans = m

	PollingIntervalSeconds = 0
	err := WatchPlan(cfg.Ctx, cfg.Client, "plan-123")
	assert.Nil(t, err)
}
