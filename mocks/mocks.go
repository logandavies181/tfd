package mocks

import (
	"os"

	"github.com/logandavies181/tfd/v2/cmd/config"

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

type Applies interface {
	tfe.Applies
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

func WithMockedFile(f *os.File, work func(f *os.File)) {
	oldFile := new(os.File)
	*oldFile = *f

	mockFile, err := os.CreateTemp("", "")
	if err != nil {
		panic(err)
	}
	*f = *mockFile

	defer func() {
		*f = *oldFile

		err := os.Remove(mockFile.Name())
		if err != nil {
			panic(err)
		}
	}()

	work(mockFile)
}
