package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultLimit  = 20
	DefaultOffset = 0
)

func ParsePagination(r *http.Request) (limit, offset int) {
	limit = DefaultLimit
	offset = DefaultOffset

	if v := r.URL.Query().Get("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	return limit, offset
}
