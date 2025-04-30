package store

import (
	"net/http"
	"strconv"
)

type PaginatedFeedQuery struct {
	Limit  int    `json:"limit" validation:"gte=1, lte=20"`
	Offset int    `json:"offset" validation:"gte=0"`
	Sort   string `json:"sort" validation:"oneof=asc desc"`
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return PaginatedFeedQuery{}, err
		}
		fq.Limit = l
	}
	offset := qs.Get("offset")
	if offset != "" {
		o, err := strconv.Atoi(offset)
		if err != nil {
			return PaginatedFeedQuery{}, err
		}
		fq.Offset = o
	}
	sort := qs.Get("sort")
	if sort != "" {

		fq.Sort = sort
	}
	return fq, nil

}
