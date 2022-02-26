package paginate

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name          string
		limit, page   int
		count         int64
		checkResponse func(t *testing.T, result *Paginate)
	}{
		{
			name:  "calculate total page correctly",
			limit: 15,
			page:  3,
			count: 40,
			checkResponse: func(t *testing.T, result *Paginate) {
				require.Equal(t, 3, result.TotalPage)
				require.Equal(t, int64(40), result.Count)
				require.Equal(t, 15, result.Limit)
				require.Equal(t, 3, result.Page)
			},
		},
		{
			name:  "next page upper bound at totalPage",
			limit: 10,
			page:  3,
			count: 30,
			checkResponse: func(t *testing.T, result *Paginate) {
				require.Equal(t, 2, result.PrevPage)
				require.Equal(t, 3, result.NextPage)
				require.Equal(t, 3, result.Page)
			},
		},
		{
			name:  "previous page lower bound at 1",
			limit: 10,
			page:  1,
			count: 30,
			checkResponse: func(t *testing.T, result *Paginate) {
				require.Equal(t, 1, result.Page)
				require.Equal(t, 1, result.PrevPage)
				require.Equal(t, 2, result.NextPage)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := New(tc.limit, tc.page, tc.count)
			tc.checkResponse(t, result)
		})
	}
}
