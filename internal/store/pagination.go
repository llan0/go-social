package store

import (
	"net/http"
	"strconv"
	"strings"
)

type PaginatedFeedQuery struct {
	Limit  int      `json:"limit" validate:"gte=1,lte=20"`
	Offset int      `json:"offset" validate:"gte=0"`
	Sort   string   `json:"sort" validate:"oneof=asc desc"`
	Tags   []string `json:"tags" validate:"max=5"`
	Search string   `json:"search" validate:"max=100"`
	//TODO: filter by time
}

func (fq PaginatedFeedQuery) Parse(r *http.Request) (PaginatedFeedQuery, error) {
	//	vl/feed?limit=20&offset=10
	qs := r.URL.Query()

	limit := qs.Get("limit")
	if limit != "" {
		l, err := strconv.Atoi(limit)
		if err != nil {
			return fq, nil
		}
		fq.Limit = l
	}

	offset := qs.Get("offset")
	if offset != "" {
		l, err := strconv.Atoi(offset)
		if err != nil {
			return fq, nil
		}
		fq.Offset = l
	}

	sort := qs.Get("sort")
	if sort != "" {
		fq.Sort = sort
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	// TODO: Since and Until
	return fq, nil
}
