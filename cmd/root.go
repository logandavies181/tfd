package cmd

import (
	"fmt"
	"os"

	"github.com/logandavies181/tfd/cmd/run"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(run.RunCmd)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tfd.yaml)")

	rootCmd.PersistentFlags().StringP("org", "o", "", "Terraform Organization to execute against")
	rootCmd.PersistentFlags().StringP("token", "t", "", "Token to use to authenticate to Terraform Cloud")

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
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".tfd.yaml" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".tfd")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvPrefix("TFD")
}
