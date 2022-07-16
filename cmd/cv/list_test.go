package cv

import (
	"os"
	"testing"

	"github.com/logandavies181/tfd/mocks"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestCvList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ws := mocks.NewMockWorkspaces(ctrl)
	cvs := mocks.NewMockConfigurationVersions(ctrl)

	gomock.InOrder(
		ws.EXPECT().
			Read(
				gomock.Any(),
				gomock.Any(),
				gomock.Eq("test-ws"),
			).
			Return(
				&tfe.Workspace{
					ID: "test-ws-id",
				},
				nil,
			),
		cvs.EXPECT().
			List(
				gomock.Any(),
				gomock.Eq("test-ws-id"),
				gomock.Nil(),
			).
			Return(
				&tfe.ConfigurationVersionList{
					Pagination: &tfe.Pagination{
						CurrentPage: 0,
						NextPage: 1,
						TotalPages: 2,
					},
					Items: []*tfe.ConfigurationVersion{
						{
							ID: "one",
						},
						{
							ID: "two",
						},
					},
				},
				nil,
			),
		cvs.EXPECT().
			List(
				gomock.Any(),
				gomock.Eq("test-ws-id"),
				gomock.Nil(),
			).
			Return(
				&tfe.ConfigurationVersionList{
					Pagination: &tfe.Pagination{
						CurrentPage: 1,
						NextPage: 1,
						TotalPages: 1,
					},
					Items: []*tfe.ConfigurationVersion{
						{
							ID: "three",
						},
						{
							ID: "four",
						},
					},
				},
				nil,
			),
	)

	cfg := mocks.MockConfig()
	cfg.Client.Workspaces = ws
	cfg.Client.ConfigurationVersions = cvs

	cvlCfg := cvListConfig{
		Config: cfg,

		MaxItems: 3,
		Workspace: "test-ws",
	}

	var cvlListErr error
	var output []byte
	mocks.WithMockedFile(os.Stdout, func(f *os.File) {
		cvlListErr = cvList(cvlCfg)

		var err error
		output, err = os.ReadFile(f.Name())
		if err != nil {
			panic(err)
		}
	})

	assert.Nil(t, cvlListErr)

	expectedStdout := "one\ntwo\nthree\n"
	assert.Equal(t, expectedStdout, string(output))
}
