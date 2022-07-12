package cmd

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	"github.com/logandavies181/tfd/pkg/git"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadConfigCmd = &cobra.Command{
	Use:          "upload-config",
	Aliases:      []string{"uc"},
	Short:        "Upload local Terraform files to Terraform Cloud",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := &uploadConfigConfig{
			Config: baseConfig,

			Path:      viper.GetString("path"),
			Workspace: viper.GetString("workspace"),
		}

		return uploadConfig(config)
	},
}

func init() {
	rootCmd.AddCommand(uploadConfigCmd)

	flags.AddPathFlag(uploadConfigCmd)
	flags.AddWorkspaceFlag(uploadConfigCmd)
}

type uploadConfigConfig struct {
	config.Config

	Path      string
	Workspace string
}

func uploadConfig(cfg *uploadConfigConfig) error {
	// get workspace id
	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	cv, err := cfg.Client.ConfigurationVersions.Create(
		cfg.Ctx,
		workspace.ID,
		tfe.ConfigurationVersionCreateOptions{})
	if err != nil {
		return err
	}

	pathToRoot, _, err := git.GetRootOfRepo(cfg.Path)
	if err != nil {
		return err
	}

	err = cfg.Client.ConfigurationVersions.Upload(cfg.Ctx, cv.UploadURL, pathToRoot)
	if err != nil {
		return err
	}

	fmt.Println("Created configuration version:", cv.ID)

	return nil
}
