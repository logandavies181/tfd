package flags

import (
	"github.com/logandavies181/tfd/cmd/config"
)

var flagValidations []func() error

func validateFlags() error {
	for _, f := range flagValidations {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}

func InitializeCmd() (*config.Config, error) {
	conf, err := config.New()
	if err != nil {
		return nil, err
	}

	err = validateFlags()
	if err != nil {
		return nil, err
	}

	return conf, nil
}
