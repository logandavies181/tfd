package workspace

import (
	"sort"

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
