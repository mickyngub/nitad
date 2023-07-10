package paginate

import "math"

type Paginate struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	PrevPage  int   `json:"prevPage"`
	NextPage  int   `json:"nextPage"`
	Count     int64 `json:"count"`
	TotalPage int   `json:"totalPage"`
}

// New create PaginationResult according to input params
func New(limit int, page int, count int64) *Paginate {
	totalPage := int(math.Ceil(float64(count) / float64(limit)))
	var nextPage int

	if page == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page + 1
	}

	var prevPage int
	if page == 1 {
		prevPage = 1
	} else {
		prevPage = page - 1
	}

	return &Paginate{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  prevPage,
		TotalPage: totalPage,
		NextPage:  nextPage,
	}
}
