package pagination

import (
	"fmt"

	"github.com/hashicorp/go-tfe"
)

// WithPagination handles initial setup and looping for a function call that may have a paginated response from the
// TFE/TFC API. The caller is expected to set the *pg to be the value from the paginated response. The caller can exit
// early by returning (true, nil)
func WithPagination(work func(pg *tfe.Pagination) (bool, error)) error {
	pg := &tfe.Pagination{
		CurrentPage: 0,
		NextPage:    1,
		TotalPages:  -1,
	}
	prevPg := tfe.Pagination{
		CurrentPage: 0,
		NextPage:    1,
		TotalPages:  -1,
	}
	for {
		if pg == nil || pg.CurrentPage == pg.TotalPages {
			break
		}
		fin, err := work(pg)
		if err != nil {
			return err
		}
		if fin {
			return nil
		}
		if *pg == prevPg {
			// indicates pg hasn't been modified as expected by work()
			return fmt.Errorf("Pagination has not progressed, exiting potential infinite loop")
		}
		prevPg = *pg
	}

	return nil
}
