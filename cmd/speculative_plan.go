package cmd

import (
	"fmt"
	"os"

	"github.com/logandavies181/tfd/cmd/config"
	"github.com/logandavies181/tfd/cmd/flags"
	"github.com/logandavies181/tfd/cmd/git"
	"github.com/logandavies181/tfd/cmd/plan"

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

		config := &speculativePlanConfig{
			Config: baseConfig,

			Path:      viper.GetString("path"),
			Workspace: viper.GetString("workspace"),
		}

		return speculativePlan(config)
	},
}

func init() {
	rootCmd.AddCommand(speculativePlanCmd)

	flags.AddPathFlag(speculativePlanCmd)
	flags.AddWorkspaceFlag(speculativePlanCmd)
}

type speculativePlanConfig struct {
	*config.Config

	Path      string
	Workspace string
}

func speculativePlan(cfg *speculativePlanConfig) error {
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
	pathToRoot, workingDir, err := git.GetRootOfRepo(cfg.Path)
	if err != nil {
		return err
	}

	if workspace.WorkingDirectory != workingDir {
		fmt.Fprintf(os.Stderr,
			"WARNING: workspace: %s will run plan using working directory: %s intead of %s (supplied).\n"+
				"Due to a limitation on the Terraform Cloud/Enterprise API, the Working Directory cannot be "+
				"overridden for a single run.\n",
			cfg.Workspace,
			workspace.WorkingDirectory,
			workingDir)
	}

	err = cfg.Client.ConfigurationVersions.Upload(cfg.Ctx, cv.UploadURL, pathToRoot)
	if err != nil {
		return err
	}

	fmt.Println(cv.ID)

	r, err := cfg.Client.Runs.Create(cfg.Ctx, tfe.RunCreateOptions{
		Workspace:            workspace,
		ConfigurationVersion: cv,
	})
	if err != nil {
		return err
	}

	fmt.Println(r.Plan.ID)

	planError := plan.WatchPlan(cfg.Ctx, cfg.Client, r.Plan.ID)
	if planError != nil {
		err, ok := planError.(plan.PlanError)
		if !ok {
			return err
		}

		fmt.Println(err)
	}

	runPlan, err := cfg.Client.Plans.Read(cfg.Ctx, r.Plan.ID)
	if err != nil {
		return err
	}

	fmt.Println("Logs available here:", runPlan.LogReadURL)

	if planError == nil {
		fmt.Println(plan.FormatResourceChanges(runPlan))
	} else {
		fmt.Println(planError)
	}

	return nil
}
