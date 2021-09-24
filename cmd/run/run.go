package run

import (
	"context"
	"fmt"
	"net/url"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
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
