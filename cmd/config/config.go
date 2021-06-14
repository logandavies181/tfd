package config

import (
	"context"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/viper"
)

type GlobalConfig struct {
	Org   string
	Token string

	Client *tfe.Client
	Ctx    context.Context
}

func GetGlobalConfig() (*GlobalConfig, error) {
	var cfg GlobalConfig

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	clientConfig := tfe.DefaultConfig()
	// TODO: add more logic for token locations
	if !(clientConfig.Token != "" && cfg.Token == "") {
		clientConfig.Token = cfg.Token // Defers to TFE_TOKEN if set
	}

	client, err := tfe.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	cfg.Client = client
	cfg.Ctx = context.TODO()

	return &cfg, nil
}

func validateGlobalConfig(globalConfig *GlobalConfig) error {
	// TODO:
	return nil
}
