package data

import (
	"hukaapp/todo/internal/validator"
	"math"
	"strings"
)

type Filter struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func (f Filter) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filter) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filter) limit() int {
	return f.PageSize
}

func (f Filter) offset() int {
	return (f.Page - 1) * f.PageSize
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func calculateMetadata(totalRecords int, page int, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func ValidateFilter(v *validator.Validator, f Filter) {
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be lower than 10 million")
	v.Check(f.PageSize > 0, "page-size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page-size", "must be lower than 100")
	v.Check(validator.In(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
