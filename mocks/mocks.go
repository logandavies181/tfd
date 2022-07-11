package mocks

import (
	"github.com/logandavies181/tfd/cmd/config"

	"github.com/hashicorp/go-tfe"
)

func MockConfig() config.Config {
	return config.Config{
		Address: "https://example.com",
		Org:     "test",
		Token:   "secret",

		Client: MockClient(),
	}
}

func MockClient() *tfe.Client {
	return &tfe.Client{
		Workspaces: &MockWorkspaces{},
	}
}

type ConfigurationVersions interface {
	tfe.ConfigurationVersions
}

type Plans interface {
	tfe.Plans
}

type Runs interface {
	tfe.Runs
}

type Workspaces interface {
	tfe.Workspaces
}
