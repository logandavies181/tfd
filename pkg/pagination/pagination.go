package pagination

import (
	"github.com/hashicorp/go-tfe"
)

func WithPagination(work func(pg *tfe.Pagination) (bool, error)) error {
	pg := &tfe.Pagination{
		NextPage:   1,
		TotalPages: -1,
	}
	for {
		if &pg == nil || pg.CurrentPage == pg.TotalPages {
			break
		}
		fin, err := work(pg)
		if err != nil {
			return err
		}
		if fin {
			return nil
		}
	}

	return nil
}
