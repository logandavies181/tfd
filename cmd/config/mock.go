package config

import (
	"github.com/hashicorp/go-tfe"
)

// Manually mock out NewClient as it makes a call to the API to get metadata

var mockNewClient bool

type NewClientCreator interface {
	NewClient(*tfe.Config) (*tfe.Client, error)
}

type newClientCreator struct{}

func (newClientCreator) NewClient(cfg *tfe.Config) (*tfe.Client, error) {
	if mockNewClient {
		return nil, nil
	}

	return tfe.NewClient(cfg)
}
