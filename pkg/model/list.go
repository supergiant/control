package model

import "strconv"

type List interface {
	QueryValues() map[string][]string
	Set(total, limit, offset int64)
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

func (l *BaseList) QueryValues() map[string][]string {
	qv := map[string][]string{
		"offset": {strconv.FormatInt(l.Offset, 10)},
		"limit":  {strconv.FormatInt(l.Limit, 10)},
	}
	for key, values := range l.Filters {
		qv["filter."+key] = values
	}
	return qv
}

func (l *BaseList) Set(total, limit, offset int64) {
	if total > 0 {
		l.Total = total
	}
	if limit > 0 {
		l.Limit = limit
	}
	if offset > 0 {
		l.Offset = offset
	}
}
