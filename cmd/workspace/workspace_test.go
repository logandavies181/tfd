package workspace

import (
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestSortWorkspacesByName(t *testing.T) {
	sortedWorkspaces := []*tfe.Workspace{
		{Name: "alice"},
		{Name: "bob"},
		{Name: "charlie"},
	}

	// permutation 1
	workspaces := []*tfe.Workspace{
		{Name: "bob"},
		{Name: "alice"},
		{Name: "charlie"},
	}
	SortWorkspacesByName(workspaces)
	assert.Equal(
		t,
		sortedWorkspaces,
		workspaces)

	// permutation 2
	workspaces = []*tfe.Workspace{
		{Name: "charlie"},
		{Name: "bob"},
		{Name: "alice"},
	}
	SortWorkspacesByName(workspaces)
	assert.Equal(
		t,
		sortedWorkspaces,
		workspaces)
}
