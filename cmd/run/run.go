package run

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
)

const (
	DESTROY = iota
	CREATE
)

var RunCmd = &cobra.Command{
	Use:           "run",
	Aliases:       []string{"r"},
	Short:         "Commands for interacting with Runs",
	SilenceErrors: true,
}

func getCurrentRun(ctx context.Context, client *tfe.Client, org, workspaceName string) (string, error) {
	workspace, err := client.Workspaces.Read(ctx, org, workspaceName)
	if err != nil {
		return "", err
	}

	if workspace.CurrentRun != nil {
		return workspace.CurrentRun.ID, nil
	} else {
		return "", fmt.Errorf("Workspace %s has no current run", workspaceName)
	}
}

func formatResourceChanges(a *tfe.Apply) string {
	return fmt.Sprintf(
		"Apply complete! Resources: %d added, %d changed, %d destroyed.",
		a.ResourceAdditions,
		a.ResourceChanges,
		a.ResourceDestructions)
}

func FormatRunUrl(address, org, workspace, runId string) (string, error) {
	addressUrl, err := url.Parse(address)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s/app/%s/workspaces/%s/runs/%s",
		addressUrl.Scheme,
		addressUrl.Host,
		org,
		workspace,
		runId), nil
}

func (cfg runStartConfig) startRun(runType int) error {
	workspace, err := cfg.Client.Workspaces.Read(cfg.Ctx, cfg.Org, cfg.Workspace)
	if err != nil {
		return err
	}

	var isDestroy bool
	switch runType {
	case DESTROY:
		isDestroy = true
	case CREATE:
		isDestroy = false
	default:
		return fmt.Errorf("Run type must be run.destroy or run.create. Unknown enum: %v", runType)
	}

	var cv *tfe.ConfigurationVersion
	if cfg.ConfigurationVersion != "" {
		cv, err = cfg.Client.ConfigurationVersions.Read(cfg.Ctx, cfg.ConfigurationVersion)
		if err != nil {
			return err
		}
	}

	var vars []*tfe.RunVariable
	for k, v := range cfg.Vars {
		vars = append(vars, &tfe.RunVariable{
			Key:   k,
			Value: v,
		})
	}

	r, err := cfg.Client.Runs.Create(
		cfg.Ctx,
		tfe.RunCreateOptions{
			AutoApply:            &cfg.FireAndForget,
			ConfigurationVersion: cv,
			IsDestroy:            &isDestroy,
			Message:              &cfg.Message,
			Refresh:              &cfg.Refresh,
			RefreshOnly:          &cfg.RefreshOnly,
			ReplaceAddrs:         cfg.Replace,
			TargetAddrs:          cfg.Targets,
			Variables:            vars,
			Workspace:            workspace,
		},
	)
	if err != nil {
		return err
	}

	fmt.Println("Started run:", r.ID)

	if cfg.Watch || cfg.AutoApply {
		err = watchAndAutoApplyRun(cfg.Ctx, cfg.Client, cfg.Org, workspace.Name, r, cfg.AutoApply, cfg.Address)
		if err != nil {
			return err
		}
	} else {
		runUrl, err := FormatRunUrl(cfg.Address, cfg.Org, workspace.Name, r.ID)
		if err != nil {
			return err
		}
		fmt.Println("View the plan in the UI:", runUrl)

	}

	return nil
}
