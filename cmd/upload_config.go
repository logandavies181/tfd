package cmd

import (
	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/git"
	"github.com/logandavies181/tfd/cmd/workspace"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var uploadConfigCmd = &cobra.Command{
	Use:          "upload-config",
	Aliases:      []string{"uc"},
	Short:        "Upload local Terraform files to Terraform Cloud",
	SilenceUsage: true,
	RunE:         uploadConfig,
}

func init() {
	rootCmd.AddCommand(uploadConfigCmd)

	uploadConfigCmd.Flags().StringP("path", "p", "", "Path to Terraform Directory")
	uploadConfigCmd.Flags().StringP("workspace", "w", "", "Terraform Cloud workspace to upload to")
}

type uploadConfigConfig struct {
	*config.GlobalConfig

	Path      string
	Workspace string
}

func getApiRunConfig(cmd *cobra.Command) (*uploadConfigConfig, error) {
	viper.BindPFlags(cmd.Flags())

	gCfg, err := config.GetGlobalConfig()
	if err != nil {
		return nil, err
	}

	var lCfg uploadConfigConfig
	err = viper.Unmarshal(&lCfg)
	if err != nil {
		return nil, err
	}

	lCfg.GlobalConfig = gCfg

	return &lCfg, nil
}

func uploadConfig(cmd *cobra.Command, _ []string) error {
	cfg, err := getApiRunConfig(cmd)
	if err != nil {
		return err
	}

	workspace, err := workspace.GetWorkspaceByName(*cfg.Client, cfg.Ctx, cfg.Org, cfg.Workspace)
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

	pathToRoot, err := git.GetRootOfRepo(cfg.Path)
	if err != nil {
		return err
	}

	err = cfg.Client.ConfigurationVersions.Upload(cfg.Ctx, cv.UploadURL, pathToRoot)
	if err != nil {
		return err
	}

	return nil
}
