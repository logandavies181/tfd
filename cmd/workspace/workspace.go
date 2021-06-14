package workspace

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/go-tfe"
)

type WorkspaceAlphabeticalSorter []*tfe.Workspace

func (was WorkspaceAlphabeticalSorter) Len() int {
	return len(was)
}

func (was WorkspaceAlphabeticalSorter) Less(i, j int) bool {
	return was[i].Name < was[j].Name
}

func (was WorkspaceAlphabeticalSorter) Swap(i, j int) {
	was[i], was[j] = was[j], was[i]
}

func SortWorkspacesByName(workspaces []*tfe.Workspace) {
	sort.Sort(WorkspaceAlphabeticalSorter(workspaces))
}

func GetWorkspaceByName(client tfe.Client, ctx context.Context, org, name string) (*tfe.Workspace, error) {
	workspaceList, err := client.Workspaces.List(ctx, org, tfe.WorkspaceListOptions{})
	if err != nil {
		return nil, err
	}

	for _, workspace := range workspaceList.Items {
		if strings.TrimSpace(workspace.Name) == strings.TrimSpace(name) {
			return workspace, nil
		}
	}

	return nil, fmt.Errorf("Could not find workspace %s", name)
}
