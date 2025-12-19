package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/marcoantoniobarcelloslimafilho/go-learning-management-system/src/internal/adapter/http/dto"
)

type paginationParams struct {
	page     int
	pageSize int
	limit    int
	offset   int
}

func parsePagination(r *http.Request) paginationParams {
	const (
		defaultPage     = 1
		defaultPageSize = 10
		maxPageSize     = 100
	)

	page := defaultPage
	pageSize := defaultPageSize

	if p := r.URL.Query().Get("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			page = v
		}
	}

	if ps := r.URL.Query().Get("page_size"); ps != "" {
		if v, err := strconv.Atoi(ps); err == nil && v > 0 {
			if v > maxPageSize {
				v = maxPageSize
			}
			pageSize = v
		}
	}

	limit := pageSize
	offset := (page - 1) * pageSize

	return paginationParams{page: page, pageSize: pageSize, limit: limit, offset: offset}
}

func buildPagination(page, pageSize int, total int64) dto.Pagination {
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))
	if totalPages == 0 {
		totalPages = 1
	}
	return dto.Pagination{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

type paginatedResponse struct {
	Data       interface{}    `json:"data"`
	Pagination dto.Pagination `json:"pagination"`
}
