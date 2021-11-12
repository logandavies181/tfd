package config

import (
	"context"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	viper.Set("address", "http://example.com/v1")
	viper.Set("token", "faketoken")
	viper.Set("org", "exampleorg")
	mockNewClient = true

	cfg, err := New()
	assert.Nil(t, err)

	assert.Equal(t, "http://example.com/v1", cfg.Address)
	assert.Equal(t, "faketoken", cfg.Token)
	assert.Equal(t, "exampleorg", cfg.Org)

	assert.Nil(t, cfg.Client) // nil because it's stubbed here
	assert.Equal(t, context.TODO(), cfg.Ctx)
}

func TestValidateConfig(t *testing.T) {
	// success case
	goodCfg := &Config{
		Token: "faketoken",
		Org: "exampleorg",
	}
	err := validateConfig(goodCfg)
	assert.Nil(t, err)

	// token failure case
	tokenFailCfg := &Config{
		Org: "exampleorg",
	}
	err = validateConfig(tokenFailCfg)
	assert.NotNil(t, err)

	// org failure case
	orgFailCfg := &Config{
		Token: "faketoken",
	}
	err = validateConfig(orgFailCfg)
	assert.NotNil(t, err)
}
