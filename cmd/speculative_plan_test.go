package cmd

import (
	"testing"

	"github.com/logandavies181/tfd/v2/mocks"

	"go.uber.org/mock/gomock"
	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestSpeculativePlan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	wsMock := mocks.NewMockWorkspaces(ctrl)
	cvMock := mocks.NewMockConfigurationVersions(ctrl)
	runsMock := mocks.NewMockRuns(ctrl)
	plansMock := mocks.NewMockPlans(ctrl)
	tru := true // go-tfe smh

	cv := tfe.ConfigurationVersion{
		ID:        "cv-1234",
		UploadURL: "http://foobar.example.com",
	}
	ws := tfe.Workspace{
		ID:               "ws-1234",
		WorkingDirectory: "workingDir",
	}
	plan := tfe.Plan{
		ID:                   "plan-1234",
		Status:               tfe.PlanFinished,
		ResourceAdditions:    1,
		ResourceChanges:      2,
		ResourceDestructions: 3,
	}

	gomock.InOrder(
		wsMock.EXPECT().
			Read(
				gomock.Any(),
				"test",
				"testWS").
			Return(
				&ws,
				nil,
			),
		cvMock.EXPECT().
			Create(
				gomock.Any(),
				"ws-1234",
				tfe.ConfigurationVersionCreateOptions{
					Speculative: &tru,
				},
			).
			Return(
				&cv,
				nil,
			),
		cvMock.EXPECT().
			Upload(
				gomock.Any(),
				"http://foobar.example.com",
				"pathToRoot", // manually mocked inside speculativePlan function
			).
			Return(nil),
		runsMock.EXPECT().
			Create(
				gomock.Any(),
				tfe.RunCreateOptions{
					Workspace:            &ws,
					ConfigurationVersion: &cv,
				},
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
			Read(gomock.Any(), gomock.Any()).
			Return(
				&plan,
				nil,
			),
		plansMock.EXPECT().
			Read(gomock.Any(), gomock.Any()).
			Return(
				&plan,
				nil,
			),
	)

	cfg := mocks.MockConfig()
	cfg.Client.Workspaces = wsMock
	cfg.Client.ConfigurationVersions = cvMock
	cfg.Client.Runs = runsMock
	cfg.Client.Plans = plansMock

	speculativePlanConfig := speculativePlanConfig{
		Config:    cfg,
		Path:      "",
		Workspace: "testWS",
		mockGit:   true,
	}

	err := speculativePlan(speculativePlanConfig)
	assert.Nil(t, err)
}
