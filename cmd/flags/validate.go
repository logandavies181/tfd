package flags

import (
	"github.com/logandavies181/tfd/v2/cmd/config"

	"github.com/spf13/cobra"
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

func InitializeCmd(cmd *cobra.Command) (conf config.Config, err error) {
	conf, err = config.New()
	if err != nil {
		return
	}

	err = validateFlags(cmd.Name())
	if err != nil {
		return
	}

	return conf, nil
}
