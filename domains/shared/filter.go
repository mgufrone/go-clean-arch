package shared

type SortDirection string
const (
	SortDesc SortDirection = "DESC"
	SortAsc = "ASC"
)
type OffsetPagination struct {
	PerPage uint64
	Page uint64
	Total uint64
}
type CursorPagination struct {
	PerPage uint64
	Next string
	Total uint64
}
type Sort struct {
	Field string
	Direction SortDirection
}
type CommonLimiter struct {
	Offset *OffsetPagination
	Cursor *CursorPagination
	Sort *Sort
	GroupBy string
	Fields []string
}

func (c *CommonLimiter) SetTotal(total uint64) *CommonLimiter {
	if c.Cursor != nil {
		c.Cursor.Total = total
	}
	if c.Offset != nil {
		c.Offset.Total = total
	}
	return c
}
func (c *CommonLimiter) Total() uint64 {
	if c.Cursor != nil {
		return c.Cursor.Total
	}
	if c.Offset != nil {
		return c.Offset.Total
	}
	return 0
}
func (c *CommonLimiter) GetPerPage() uint64 {
	if c.Cursor != nil {
		return c.Cursor.PerPage
	}
	if c.Offset != nil {
		return c.Offset.PerPage
	}
	return 0
}

func DefaultSort() *Sort {
	return &Sort{
		Field:     "created_at",
		Direction: SortDesc,
	}
}
func DefaultLimiter(limit uint64) *CommonLimiter  {
	if limit == 0 {
		limit = 10
	}
	return &CommonLimiter{
		Offset: &OffsetPagination{
			PerPage: limit,
			Page:    1,
		},
		Sort: DefaultSort(),
	}
}