package flags

import (
	"fmt"

	"github.com/logandavies181/tfd/v2/cmd/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagValidations = make(map[string][]func() error)
)

func validateFlags(name string) error {
	for _, f := range flagValidations[name] {
		err := f()
		if err != nil {
			return err
		}
	}

	return nil
}

func addValidation(name string, validationFunc func() error) {
	if v, ok := flagValidations[name]; ok {
		flagValidations[name] = append(v, validationFunc)
	} else {
		flagValidations[name] = []func() error{validationFunc}
	}
}

func InitializeCmd(cmd *cobra.Command) (config.Config, error) {
	conf, err := config.New()
	if err != nil {
		return config.Config{}, fmt.Errorf("could not create new config object: %w", err)
	}

	err = viper.BindPFlags(cmd.Flags())
	if err != nil {
		return config.Config{}, fmt.Errorf("could not bind flags: %w", err)
	}

	err = validateFlags(cmd.Name())
	if err != nil {
		return config.Config{}, err // just let the validation error display as-is
	}

	return conf, nil
}
