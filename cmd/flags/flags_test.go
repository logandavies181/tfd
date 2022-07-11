package flags

import (
	"errors"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const testCommandName = "test123123"

func TestUnvalidatedFlags(t *testing.T) {
	cmd := &cobra.Command{
		Use: testCommandName,
	}
	AddPathFlag(cmd)
	AddAutoApplyFlag(cmd)
	AddMessageFlag(cmd)
	AddWatchFlag(cmd)
	AddRefreshFlag(cmd)
	AddRefreshOnlyFlag(cmd)
	AddReplaceFlag(cmd)
	AddTargetsFlag(cmd)
	AddRunIdFlag(cmd)

	err := validateFlags(testCommandName)
	assert.Nil(t, err)
}

func TestAddWorkSpaceFlag(t *testing.T) {
	cmd := &cobra.Command{
		Use: testCommandName,
	}
	AddWorkspaceFlag(cmd)
	viper.Set("workspace", "foobar")

	// expect success
	err := validateFlags(testCommandName)
	assert.Nil(t, err)

	// expect failure
	viper.Set("workspace", "")
	err = validateFlags(testCommandName)
	assert.NotNil(t, err)
}

func TestValidateFlags(t *testing.T) {
	successCommand := func() error {
		return nil
	}
	failCOmmand := func() error {
		return errors.New("")
	}

	// expect success
	flagValidations[testCommandName] = []func() error{
		successCommand,
	}
	err := validateFlags(testCommandName)
	assert.Nil(t, err)

	// expect failure
	flagValidations[testCommandName] = []func() error{
		failCOmmand,
	}
	err = validateFlags(testCommandName)
	assert.NotNil(t, err)
}
