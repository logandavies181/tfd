package vars

import (
	"strings"

	"github.com/hashicorp/go-tfe"
	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var varsSetCmd = &cobra.Command{
	Use:          "set",
	Aliases:      []string{"s"},
	Short:        "Create or update a variable on a workspace",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		config := varsSetConfig{
			Config: baseConfig,

			Category:    viper.GetString("category"),
			Description: viper.GetString("description"),
			Hcl:         viper.GetBool("hcl"),
			Key:         viper.GetString("key"),
			NoClobber:   viper.GetBool("no-clobber"),
			Sensitive:   viper.GetBool("sensitive"),
			Value:       viper.GetString("value"),
			Workspace:   viper.GetString("workspace"),
		}

		return varsSet(config)
	},
}

func init() {
	VarsCmd.AddCommand(varsSetCmd)

	flags.AddCategoryFlag(varsSetCmd)
	flags.AddDescriptionFlag(varsSetCmd)
	flags.AddHclFlag(varsSetCmd)
	flags.AddKeyFlag(varsSetCmd)
	flags.AddNoClobberFlag(varsSetCmd)
	flags.AddSensitiveFlag(varsSetCmd)
	flags.AddValueFlag(varsSetCmd)
	flags.AddWorkspaceFlag(varsSetCmd)
}

type varsSetConfig struct {
	config.Config

	Category    string
	Description string
	Hcl         bool
	Key         string
	NoClobber   bool
	Sensitive   bool
	Value       string
	Workspace   string
}

func varsSet(cfg varsSetConfig) error {
	ws, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	v, err := cfg.Client.Variables.Create(cfg.Ctx, ws.ID, tfe.VariableCreateOptions{
		Category: categoryType(cfg.Category),
		Description: &cfg.Description,
		HCL: &cfg.Hcl,
		Key: &cfg.Key,
		Sensitive: &cfg.Sensitive,
		Value: &cfg.Value,
	})
	if err != nil && !strings.Contains(err.Error(), "Key has already been taken") {
		return err
	}

	if cfg.NoClobber {
		return nil
	}

	_, err = cfg.Client.Variables.Update(cfg.Ctx, ws.ID, v.ID, tfe.VariableUpdateOptions{
		Category: categoryType(cfg.Category),
		Description: &cfg.Description,
		HCL: &cfg.Hcl,
		Key: &cfg.Key,
		Sensitive: &cfg.Sensitive,
		Value: &cfg.Value,
	})

	return err
}
