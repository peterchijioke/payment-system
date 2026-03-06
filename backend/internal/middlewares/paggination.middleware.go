package middlewares

import (
	"math"

	"gorm.io/gorm"
)

type PaginationResult[T any] struct {
	TotalDocs     int64 `json:"totalDocs"`
	Limit         int   `json:"limit"`
	TotalPages    int   `json:"totalPages"`
	Page          int   `json:"page"`
	PagingCounter int   `json:"pagingCounter"`
	HasPrevPage   bool  `json:"hasPrevPage"`
	HasNextPage   bool  `json:"hasNextPage"`
	PrevPage      int   `json:"prevPage"`
	NextPage      int   `json:"nextPage"`
	HasMore       bool  `json:"hasMore"`
	Docs          []T   `json:"docs"`
}

func PaginateWithMeta[T any](db *gorm.DB, page int, limit int, model *[]T) (*PaginationResult[T], error) {
	if limit > 100 {
		limit = 100
	} else if limit <= 0 {
		limit = 10
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	var total int64
	db.Model(model).Count(&total)

	err := db.Offset(offset).Limit(limit).Find(model).Error
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasPrev := page > 1
	hasNext := page < totalPages

	result := &PaginationResult[T]{
		TotalDocs:     total,
		Limit:         limit,
		TotalPages:    totalPages,
		Page:          page,
		PagingCounter: offset + 1,
		HasPrevPage:   hasPrev,
		HasNextPage:   hasNext,
		PrevPage: func() int {
			if hasPrev {
				return page - 1
			}
			return -1
		}(),
		NextPage: func() int {
			if hasNext {
				return page + 1
			}
			return -1
		}(),
		HasMore: offset+limit < int(total),
		Docs:    *model,
	}

	return result, nil
}
