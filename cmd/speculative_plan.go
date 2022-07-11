package cmd

import (
	"fmt"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	"github.com/logandavies181/tfd/cmd/plan"
	"github.com/logandavies181/tfd/cmd/run"
	"github.com/logandavies181/tfd/pkg/git"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var speculativePlanCmd = &cobra.Command{
	Use:          "speculative-plan",
	Aliases:      []string{"spec", "sp"},
	Short:        "Run a speculative plan using local files. Also works with VCS-integrated Workspaces",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		baseConfig, err := flags.InitializeCmd(cmd)
		if err != nil {
			return err
		}

		rsc := run.RunStartConfig{
			AutoApply:     viper.GetBool("auto-apply"),
			FireAndForget: viper.GetBool("fire-and-forget"),
			Message:       viper.GetString("message"),
			Refresh:       viper.GetBool("refresh"),
			RefreshOnly:   viper.GetBool("refresh-only"),
			Replace:       viper.GetStringSlice("replace"),
			Targets:       viper.GetStringSlice("targets"),
			Watch:         viper.GetBool("watch"),
			Workspace:     viper.GetString("workspace"),
		}

		config := speculativePlanConfig{
			Config: baseConfig,

			Path:      viper.GetString("path"),
			Workspace: viper.GetString("workspace"),
			RunStartConfig: rsc,
		}

		return speculativePlan(config)
	},
}

func init() {
	rootCmd.AddCommand(speculativePlanCmd)

	flags.AddPathFlag(speculativePlanCmd)
	flags.AddWorkspaceFlag(speculativePlanCmd)
	flags.AddAutoApplyFlag(speculativePlanCmd)
	flags.AddFireAndForgetFlag(speculativePlanCmd)
	flags.AddMessageFlag(speculativePlanCmd)
	flags.AddRefreshFlag(speculativePlanCmd)
	flags.AddRefreshOnlyFlag(speculativePlanCmd)
	flags.AddReplaceFlag(speculativePlanCmd)
	flags.AddTargetsFlag(speculativePlanCmd)
	flags.AddWatchFlag(speculativePlanCmd)
}

type speculativePlanConfig struct {
	config.Config

	Path      string
	Workspace string
	RunStartConfig       run.RunStartConfig

	mockGit bool
}

func speculativePlan(cfg speculativePlanConfig) error {
	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	tru := true
	cv, err := cfg.Client.ConfigurationVersions.Create(
		cfg.Ctx,
		workspace.ID,
		tfe.ConfigurationVersionCreateOptions{
			Speculative: &tru, // lib demands a *bool
		})
	if err != nil {
		return err
	}
	var pathToRoot string
	if cfg.mockGit {
		pathToRoot = "pathToRoot"
	} else {
		pathToRoot, _, err = git.GetRootOfRepo(cfg.Path)
		if err != nil {
			return err
		}
	}

	err = cfg.Client.ConfigurationVersions.Upload(cfg.Ctx, cv.UploadURL, pathToRoot)
	if err != nil {
		return err
	}

	fmt.Println("Created configuration version:", cv.ID)

		r, err := cfg.Client.Runs.Create(cfg.Ctx, tfe.RunCreateOptions{
		Workspace:            workspace,
		ConfigurationVersion: cv,
	})
	if err != nil {
		return err
	}

	fmt.Println(r.Plan.ID)

	runUrl, err := run.FormatRunUrl(cfg.Address, cfg.Org, cfg.Workspace, r.ID)
	if err != nil {
		return err
	}
	fmt.Println("View the run in the UI:", runUrl)

	planError := plan.WatchPlan(cfg.Ctx, cfg.Client, r.Plan.ID)
	if planError != nil {
		err, ok := planError.(plan.PlanError)
		if !ok {
			return err
		}
	}

	runPlan, err := cfg.Client.Plans.Read(cfg.Ctx, r.Plan.ID)
	if err != nil {
		return err
	}

	if planError == nil {
		fmt.Println(plan.FormatResourceChanges(runPlan))
	} else {
		fmt.Println(planError)
	}

	return nil
}
