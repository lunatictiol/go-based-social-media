package store

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validation:"gte=1, lte=20"`
	Offset int      `json:"offset" validation:"gte=0"`
	Sort   string   `json:"sort" validation:"oneof=asc desc"`
	Tags   []string `json:"tags" validation:"max=5"`
	Search string   `json:"search" validation:"max=500"`
	Since  string   `json:"since"`
	Until  string   `json:"until"`
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

	tags := qs.Get("tags")
	if tags != "" {

		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {

		fq.Search = search
	}
	since := qs.Get("since")
	if since != "" {

		fq.Since = parseTime(since)
	}
	until := qs.Get("until")
	if search != "" {

		fq.Until = parseTime(until)
	}

	return fq, nil

}

func parseTime(timeString string) string {
	t, err := time.Parse(time.DateTime, timeString)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}
