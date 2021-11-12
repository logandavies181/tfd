package cmd

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	"github.com/logandavies181/tfd/cmd/git"

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

			Path:               viper.GetString("path"),
			Workspace:          viper.GetString("workspace"),
			NoUpdateWorkingDir: viper.GetBool("no-update-workingdir"),
		}

		return uploadConfig(config)
	},
}

func init() {
	rootCmd.AddCommand(uploadConfigCmd)

	flags.AddPathFlag(uploadConfigCmd)
	flags.AddWorkspaceFlag(uploadConfigCmd)
	flags.AddNoUpdateWorkingdirFlag(uploadConfigCmd)
}

type uploadConfigConfig struct {
	*config.Config

	Path               string
	Workspace          string
	NoUpdateWorkingDir bool
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

	pathToRoot, workingDir, err := git.GetRootOfRepo(cfg.Path)
	if err != nil {
		return err
	}

	if !cfg.NoUpdateWorkingDir {
		_, err = cfg.Client.Workspaces.Update(cfg.Ctx, cfg.Org, cfg.Workspace, tfe.WorkspaceUpdateOptions{
			WorkingDirectory: &workingDir,
		})
		if err != nil {
			return fmt.Errorf("Failed to update workspace working directory: %s", err)
		}
	}

	err = cfg.Client.ConfigurationVersions.Upload(cfg.Ctx, cv.UploadURL, pathToRoot)
	if err != nil {
		return err
	}

	fmt.Println(cv.ID)

	return nil
}
