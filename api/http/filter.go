package http

import (
	"recipes/domains/shared"
	"strings"
)

type IFilter interface {
	Transform() *shared.CommonLimiter
}

type Pagination struct {
	PerPage int `form:"per_page" json:"per_page"`
	Page int `form:"page" json:"page"`
	Cursor string `form:"cursor" json:"-"`
}
type Sort struct {
	SortBy string `form:"sort_by" json:"sort_by"`
	Sort string `form:"sort" json:"sort"`
}

type CommonFilter struct {
	Pagination
	Sort
}

func (c *CommonFilter) Transform() *shared.CommonLimiter {
	var (
		offset *shared.OffsetPagination
		cursor *shared.CursorPagination
	)
	if c.PerPage == 0 {
		c.PerPage = 10
	}
	if c.Page == 0 {
		c.Page = 1
	}
	if c.PerPage > 0 || (c.PerPage == 0 && c.Cursor == ""){
		offset = &shared.OffsetPagination{
			PerPage: uint64(c.PerPage),
			Page:    uint64(c.Page),
		}
	}
	if c.Cursor != "" {
		cursor = &shared.CursorPagination{
			PerPage: uint64(c.PerPage),
			Next:    c.Cursor,
		}
	}
	dir := shared.SortDesc
	if strings.ToLower(c.SortBy) == "asc" {
		dir = shared.SortAsc
	}
	return &shared.CommonLimiter{
		Offset: offset,
		Cursor: cursor,
		Sort:   &shared.Sort{
			Field:     c.SortBy,
			Direction: dir,
		},
	}
}
