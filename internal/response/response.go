package response

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type PaginatedResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}

func NewPaginated[T any](data []T, total, limit, offset int) PaginatedResponse[T] {
	page := 1
	totalPages := 1

	if limit > 0 {
		page = offset/limit + 1
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	if totalPages < 1 {
		totalPages = 1
	}

	if data == nil {
		data = []T{}
	}

	return PaginatedResponse[T]{
		Data: data,
		Pagination: Pagination{
			Page:       page,
			PageSize:   limit,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}
}

func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error": "%s"}`, err)
}
