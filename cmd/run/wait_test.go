package run

import (
	"testing"

	"github.com/logandavies181/tfd/cmd/plan"
	"github.com/logandavies181/tfd/mocks"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestWaitForQueueStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	PollingIntervalSeconds = 0
	plan.PollingIntervalSeconds = 0

	wsMock := mocks.NewMockWorkspaces(ctrl)

	workspaceWithCurrentRun := tfe.Workspace{
		ID:               "ws-1234",
		WorkingDirectory: "workingDir",
		CurrentRun: &tfe.Run{
			ID:     "run-1234",
			Status: "not_finished",
		},
	}

	gomock.InOrder(
		wsMock.EXPECT().
			Read(
				gomock.Any(),
				"test",
				"testWS",
			).
			Return(
				&workspaceWithCurrentRun,
				nil,
			),
		wsMock.EXPECT().
			Read(
				gomock.Any(),
				"test",
				"testWS",
			).
			Return(
				&workspaceWithCurrentRun,
				nil,
			),
	)

	cfg := mocks.MockConfig()
	cfg.Client.Workspaces = wsMock

	err := waitForQueueStatus(cfg.Ctx, cfg.Client, cfg.Org, "testWS", "run-1234")
	assert.Nil(t, err)
}

func TestWaitForQueueStatusErrorNotCurrentRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	PollingIntervalSeconds = 0
	plan.PollingIntervalSeconds = 0

	wsMock := mocks.NewMockWorkspaces(ctrl)

	workspaceWithCurrentRun := tfe.Workspace{
		ID:               "ws-1234",
		WorkingDirectory: "workingDir",
		CurrentRun: &tfe.Run{
			ID:     "run-1234",
			Status: "not_finished",
		},
	}

	gomock.InOrder(
		wsMock.EXPECT().
			Read(
				gomock.Any(),
				"test",
				"testWS",
			).
			Return(
				&workspaceWithCurrentRun,
				nil,
			),
		wsMock.EXPECT().
			Read(
				gomock.Any(),
				"test",
				"testWS",
			).
			Return(
				&workspaceWithCurrentRun,
				nil,
			),
	)

	cfg := mocks.MockConfig()
	cfg.Client.Workspaces = wsMock

	err := waitForQueueStatus(cfg.Ctx, cfg.Client, cfg.Org, "testWS", "run-5678")
	assert.Error(t, err)
}

func TestWatchRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	PollingIntervalSeconds = 0
	plan.PollingIntervalSeconds = 0

	runsMock := mocks.NewMockRuns(ctrl)

	notFInishedRun := &tfe.Run{
		Status: "not_finished",
	}

	gomock.InOrder(
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				notFInishedRun,
				nil,
			),
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				notFInishedRun,
				nil,
			),
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				&tfe.Run{
					Status: tfe.RunApplied,
				},
				nil,
			),
	)

	cfg := mocks.MockConfig()
	cfg.Client.Runs = runsMock

	err := watchRun(cfg.Ctx, cfg.Client, "run-1234")
	assert.Nil(t, err)
}

func TestWatchRunErrored(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	PollingIntervalSeconds = 0
	plan.PollingIntervalSeconds = 0

	runsMock := mocks.NewMockRuns(ctrl)

	erroredRun := &tfe.Run{
		Status: tfe.RunErrored,
	}

	gomock.InOrder(
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				erroredRun,
				nil,
			),
	)

	cfg := mocks.MockConfig()
	cfg.Client.Runs = runsMock

	err := watchRun(cfg.Ctx, cfg.Client, "run-1234")
	assert.Error(t, err)
}

func TestWatchAndAutoApplyRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	PollingIntervalSeconds = 0
	plan.PollingIntervalSeconds = 0

	wsMock := mocks.NewMockWorkspaces(ctrl)
	plansMock := mocks.NewMockPlans(ctrl)
	runsMock := mocks.NewMockRuns(ctrl)
	appliesMock := mocks.NewMockApplies(ctrl)

	notFInishedRun := &tfe.Run{
		Status: "not_finished",
	}

	workspaceWithCurrentRun := tfe.Workspace{
		ID:               "ws-1234",
		WorkingDirectory: "workingDir",
		CurrentRun: &tfe.Run{
			ID:     "run-1234",
			Status: "not_finished",
		},
	}

	gomock.InOrder(
		// region waitForQueueStatus
		wsMock.EXPECT().
			Read(
				gomock.Any(),
				"test",
				"testWS",
			).
			Return(
				&workspaceWithCurrentRun,
				nil,
			),
		wsMock.EXPECT().
			Read(
				gomock.Any(),
				"test",
				"testWS",
			).
			Return(
				&workspaceWithCurrentRun,
				nil,
			),
		// endregion

		// local
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				&tfe.Run{
					Plan: &tfe.Plan{
						ID: "plan-1234",
					},
				},
				nil,
			),
		// endlocal

		// region watchPlan
		plansMock.EXPECT().
			Read(
				gomock.Any(),
				gomock.Eq("plan-1234"),
			).
			Return(
				&tfe.Plan{
					Status: "some status",
				},
				nil,
			),
		plansMock.EXPECT().
			Read(
				gomock.Any(),
				gomock.Eq("plan-1234"),
			).
			Return(
				&tfe.Plan{
					Status: tfe.PlanFinished,
				},
				nil,
			),
		// endregion watchPlan

		// local
		plansMock.EXPECT().
			Read(
				gomock.Any(),
				gomock.Eq("plan-1234"),
			).
			Return(
				&tfe.Plan{
					Status: "some status",
				},
				nil,
			),
		// endlocal

		// region waitForRunToBeConfirmable
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				&tfe.Run{
					Status:  tfe.RunPolicyChecking,
					Actions: &tfe.RunActions{},
				},
				nil,
			),
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				&tfe.Run{
					ID: "run-1234",
					Actions: &tfe.RunActions{
						IsConfirmable: true,
					},
				},
				nil,
			),
		runsMock.EXPECT().
			Apply(
				gomock.Any(),
				"run-1234",
				gomock.Any(),
			).
			Return(
				nil,
			),
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				&tfe.Run{
					ID: "run-1234",
					Actions: &tfe.RunActions{
						IsConfirmable: false,
					},
					Status: tfe.RunApplying,
				},
				nil,
			),
		// endregion waitForRunToBeConfirmable

		// region watchRun
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				notFInishedRun,
				nil,
			),
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				notFInishedRun,
				nil,
			),
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				&tfe.Run{
					Apply: &tfe.Apply{
						ID: "apply-1234",
					},
					Status: tfe.RunApplied,
				},
				nil,
			),
		// endregion watchRun

		// local
		runsMock.EXPECT().
			Read(
				gomock.Any(),
				"run-1234",
			).
			Return(
				&tfe.Run{
					ID:     "run-1234",
					Status: tfe.RunApplied,
					Apply: &tfe.Apply{
						ID: "apply-1234",
					},
				},
				nil,
			),
		appliesMock.EXPECT().
			Read(
				gomock.Any(),
				"apply-1234",
			).
			Return(
				&tfe.Apply{
					ResourceAdditions:    1,
					ResourceDestructions: 1,
					ResourceChanges:      1,
				},
				nil,
			),
		// endlocal
	)

	cfg := mocks.MockConfig()
	cfg.Client.Applies = appliesMock
	cfg.Client.Runs = runsMock
	cfg.Client.Plans = plansMock
	cfg.Client.Workspaces = wsMock

	runToWatch := &tfe.Run{
		ID: "run-1234",
	}

	err := watchAndAutoApplyRun(cfg.Ctx, cfg.Client, "test", "testWS", runToWatch, true, cfg.Address)
	assert.Nil(t, err)
}
