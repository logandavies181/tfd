package mocks

import (
	"github.com/logandavies181/tfd/cmd/config"

	"github.com/hashicorp/go-tfe"
)

func MockConfig() *config.Config {
	return &config.Config{
		Org:   "test",
		Token: "secret",

		Client: MockClient(),
	}
}

func MockClient() *tfe.Client {
	return &tfe.Client{
		Workspaces: &MockWorkspaces{},
	}
}

type Workspaces interface {
	tfe.Workspaces
}
