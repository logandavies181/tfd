package pagination

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-tfe"
	"github.com/stretchr/testify/assert"
)

func TestWithPaginationNoEdit(t *testing.T) {
	err := WithPagination(func(pg *tfe.Pagination) (bool, error) {
		return false, nil
	})
	assert.Error(t, err)
}

func TestPaginateThreeTimes(t *testing.T) {
	expectedLoops := 3
	count := 0

	err := WithPagination(func(pg *tfe.Pagination) (bool, error) {
		pg.NextPage++
		pg.CurrentPage++
		pg.TotalPages = expectedLoops

		count++

		return false, nil
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedLoops, count)
}

func TestPaginateExitEarly(t *testing.T) {
	expectedLoops := 2
	count := 0

	err := WithPagination(func(pg *tfe.Pagination) (bool, error) {
		pg.NextPage++
		pg.CurrentPage++
		pg.TotalPages = expectedLoops + 1

		count++
		if count == 2 {
			return true, nil
		}

		return false, nil
	})

	assert.Nil(t, err)
	assert.Equal(t, expectedLoops, count)
}

func TestPaginatePropagateCallerError(t *testing.T) {
	testError := fmt.Errorf("Test error")

	err := WithPagination(func(pg *tfe.Pagination) (bool, error) {
		return false, testError
	})

	assert.Equal(t, testError, err)
}
