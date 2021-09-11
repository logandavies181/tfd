package flags

import (
	"github.com/logandavies181/tfd/cmd/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func InitializeCmd(cmd *cobra.Command) (*config.Config, error) {
	viper.BindPFlags(cmd.Flags())

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
