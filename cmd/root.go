package cmd

import (
	"fmt"
	"os"

	"github.com/logandavies181/tfd/cmd/cv"
	"github.com/logandavies181/tfd/cmd/run"
	"github.com/logandavies181/tfd/cmd/vars"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	defaultTerraformCloudURI = "https://app.terraform.io/api/v2"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "tfd",
	Short:         "Command line helper for Terraform Cloud and Terraform Enterprise",
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(semver string) {
	rootCmd.Version = semver
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(cv.CvCmd)
	rootCmd.AddCommand(run.RunCmd)
	rootCmd.AddCommand(vars.VarsCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tfd.yaml)")

	rootCmd.PersistentFlags().StringP("org", "o", "", "Terraform Organization to execute against")
	rootCmd.PersistentFlags().StringP("token", "t", "", "Token to use to authenticate to Terraform Cloud")
	rootCmd.PersistentFlags().StringP("address", "", defaultTerraformCloudURI, "Full Terraform Cloud/Enterprise API URI")

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tfd.yaml" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tfd")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("TFD")

	err := viper.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// No config file, no big deal
		default:
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}
