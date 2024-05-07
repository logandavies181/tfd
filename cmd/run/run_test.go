package run

import (
	"testing"

	"github.com/logandavies181/tfd/v2/cmd/plan"
	"github.com/logandavies181/tfd/v2/mocks"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetCurrentRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	PollingIntervalSeconds = 0

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
	)

	cfg := mocks.MockConfig()
	cfg.Client.Workspaces = wsMock

	currentRunID, err := getCurrentRun(cfg.Ctx, cfg.Client, "test", "testWS")
	assert.Nil(t, err)
	assert.Equal(t, "run-1234", currentRunID)
}

func TestFormatResourceChanges(t *testing.T) {
	input := &tfe.Apply{
		ResourceAdditions:    0,
		ResourceDestructions: 1,
		ResourceChanges:      2,
	}

	output := formatResourceChanges(input)
	assert.Equal(t, "Apply complete! Resources: 0 added, 2 changed, 1 destroyed.", output)
}

func TestFormatRunUrl(t *testing.T) {
	output, err := FormatRunUrl("https://example.com", "test", "testWS", "run-1234")
	assert.Nil(t, err)
	assert.Equal(t, "https://example.com/app/test/workspaces/testWS/runs/run-1234", output)
}

func TestStartRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	PollingIntervalSeconds = 0
	plan.PollingIntervalSeconds = 0

	wsMock := mocks.NewMockWorkspaces(ctrl)
	plansMock := mocks.NewMockPlans(ctrl)
	runsMock := mocks.NewMockRuns(ctrl)
	appliesMock := mocks.NewMockApplies(ctrl)
	cvsMock := mocks.NewMockConfigurationVersions(ctrl)

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

	ws := &tfe.Workspace{
		ID:   "ws-1234",
		Name: "testWS",
	}
	cv := &tfe.ConfigurationVersion{}
	fals := false
	message := "hello, world"

	gomock.InOrder(
		// region local
		wsMock.EXPECT().
			Read(
				gomock.Any(),
				"test",
				"testWS",
			).
			Return(
				ws,
				nil,
			),
		cvsMock.EXPECT().
			Read(
				gomock.Any(),
				"cv-1234",
			).
			Return(
				cv,
				nil,
			),
		runsMock.EXPECT().
			Create(
				gomock.Any(),
				gomock.Eq(tfe.RunCreateOptions{
					AutoApply:            &fals,
					ConfigurationVersion: cv,
					IsDestroy:            &fals,
					Message:              &message,
					Refresh:              &fals,
					RefreshOnly:          &fals,
					ReplaceAddrs:         []string{},
					TargetAddrs:          []string{},
					Workspace:            ws,
				}),
			).
			Return(
				&tfe.Run{
					ID: "run-1234",
				},
				nil,
			),
		// endregion local

		// region watchAndAutoApplyRun
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
		// endregion watchAndAutoApplyRun
	)

	cfg := mocks.MockConfig()
	cfg.Client.Applies = appliesMock
	cfg.Client.Runs = runsMock
	cfg.Client.Plans = plansMock
	cfg.Client.Workspaces = wsMock
	cfg.Client.ConfigurationVersions = cvsMock

	c := &runStartConfig{
		Config:               cfg,
		AutoApply:            true,
		ConfigurationVersion: "cv-1234",
		FireAndForget:        false,
		Message:              "hello, world",
		Refresh:              false,
		RefreshOnly:          false,
		Replace:              []string{},
		Targets:              []string{},
		Vars:                 make(map[string]string),
		Watch:                false,
		Workspace:            "testWS",
	}

	err := c.startRun(CREATE)
	assert.Nil(t, err)
}
