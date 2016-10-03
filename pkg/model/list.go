package model

import "strconv"

type List interface {
	QueryValues() map[string][]string
	SetFilters(map[string][]string)
}

type BaseList struct {
	// Filters input looks like this:
	// {
	// 	"name": [
	// 		"this", "that"
	// 	],
	// 	"other_field": [
	// 		"thingy"
	// 	]
	// }
	// The values of each key are joined as OR queries.
	// Multiple keys are joined as AND queries.
	// The above translates to "(name=this OR name=that) AND (other_fied=thingy)".
	Filters map[string][]string `json:"filters"`

	// Pagination
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
	Total  int64 `json:"total"`
}

func (l BaseList) QueryValues() map[string][]string {
	qv := map[string][]string{
		"offset": []string{strconv.FormatInt(l.Offset, 10)},
		"limit":  []string{strconv.FormatInt(l.Limit, 10)},
	}
	for key, values := range l.Filters {
		qv["filter."+key] = values
	}
	return qv
}

func (l BaseList) SetFilters(filters map[string][]string) {
	l.Filters = filters
}
