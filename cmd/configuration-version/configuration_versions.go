package configuration_version

import (
	"context"

	"github.com/hashicorp/go-tfe"
)

func GetConfigurationVersionById(ctx context.Context, client *tfe.Client, workspaceId, Id string) (*tfe.ConfigurationVersion, error) {

	return nil, nil
}

func WithPagination(work func(pagination *tfe.Pagination) error, breakFunc func() bool) error {
	pagination := &tfe.Pagination{
		NextPage:   1,
		TotalPages: -1,
	}
	for {
		if pagination == nil || pagination.CurrentPage == pagination.TotalPages {
			break
		}
		err := work(pagination)
		if err != nil {
			return err
		}

		if breakFunc() {
			break
		}
	}

	return nil
}
