package cmd

import (
	"testing"

	"github.com/logandavies181/tfd/v2/mocks"

	"github.com/golang/mock/gomock"
	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestListWorkspacesOnePage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockWorkspaces(ctrl)

	m.
		EXPECT().
		List(gomock.Any(), "test", gomock.Any()).
		Return(&tfe.WorkspaceList{
			Pagination: &tfe.Pagination{
				CurrentPage: 1,
				TotalPages:  1,
			},
			Items: []*tfe.Workspace{
				{
					Name: "one",
				},
			},
		}, nil).
		Times(1)

	cfg := mocks.MockConfig()
	cfg.Client.Workspaces = m

	err := listWorkspaces(cfg)
	assert.Nil(t, err)
}

func TestListWorkspacesManyPages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockWorkspaces(ctrl)

	gomock.InOrder(
		m.EXPECT().
			List(gomock.Any(), "test", &tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 1,
				},
			}).
			Return(&tfe.WorkspaceList{
				Pagination: &tfe.Pagination{
					CurrentPage: 1,
					NextPage:    2,
					TotalPages:  3,
				},
				Items: []*tfe.Workspace{
					{
						Name: "one",
					},
				},
			}, nil),
		m.EXPECT().
			List(gomock.Any(), "test", &tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 2,
				},
			}).
			Return(&tfe.WorkspaceList{
				Pagination: &tfe.Pagination{
					CurrentPage: 2,
					NextPage:    3,
					TotalPages:  3,
				},
				Items: []*tfe.Workspace{
					{
						Name: "two",
					},
				},
			}, nil),
		m.EXPECT().
			List(gomock.Any(), "test", &tfe.WorkspaceListOptions{
				ListOptions: tfe.ListOptions{
					PageNumber: 3,
				},
			}).
			Return(&tfe.WorkspaceList{
				Pagination: &tfe.Pagination{
					CurrentPage: 3,
					TotalPages:  3,
				},
				Items: []*tfe.Workspace{
					{
						Name: "three",
					},
				},
			}, nil),
	)

	cfg := mocks.MockConfig()
	cfg.Client.Workspaces = m

	err := listWorkspaces(cfg)
	assert.Nil(t, err)
}
