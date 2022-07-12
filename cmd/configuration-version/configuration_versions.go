package configuration_version

import (
	"context"
	"fmt"

	"github.com/logandavies181/tfd/pkg/pagination"

	"github.com/hashicorp/go-tfe"
)

func GetConfigurationVersionById(ctx context.Context, client *tfe.Client, workspaceId, id string) (*tfe.ConfigurationVersion, error) {
	var cv *tfe.ConfigurationVersion
	pagination.WithPagination(func(pg *tfe.Pagination) (bool, error) {
		cvl, err := client.ConfigurationVersions.List(ctx, workspaceId, &tfe.ConfigurationVersionListOptions{
			ListOptions: tfe.ListOptions{
				PageNumber: pg.NextPage,
			},
		})
		if err != nil {
			return false, err
		}

		for _, item := range cvl.Items {
			if item.ID == id {
				cv = item

				return true, nil
			}
		}

		if cvl.Pagination != nil {
			*pg = *cvl.Pagination
		}

		return false, nil
	})

	if cv == nil {
		return nil, fmt.Errorf("Could not find configuration version Id: %v", id)
	}

	return cv, nil
}
