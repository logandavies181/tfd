package flags

import (
	"github.com/logandavies181/tfd/cmd/config"

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

func InitializeCmd(cmd *cobra.Command) (*config.Config, error) {
	viper.BindPFlags(cmd.Flags())

	conf, err := config.New()
	if err != nil {
		return nil, err
	}

	err = validateFlags(cmd.Name())
	if err != nil {
		return nil, err
	}

	return conf, nil
}
