package pagination

const (
	DefaultLimit int64 = 10
	MaxLimit     int64 = 100
	DefaultPage  int64 = 1
)

type Params struct {
	Page  int64
	Limit int64
}

func (p Params) Normalize() Params {
	page := p.Page
	if page <= 0 {
		page = DefaultPage
	}
	limit := p.Limit
	if limit <= 0 {
		limit = DefaultLimit
	} else if limit > MaxLimit {
		limit = MaxLimit
	}
	return Params{Page: page, Limit: limit}
}

func (p Params) Skip() int64 {
	return (p.Page - 1) * p.Limit
}

type Result[T any] struct {
	Items      []T   `json:"items"`
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	TotalCount int64 `json:"total_count"`
	TotalPages int64 `json:"total_pages"`
}

func NewResult[T any](items []T, p Params, totalCount int64) Result[T] {
	totalPages := totalCount / p.Limit
	if totalCount%p.Limit != 0 {
		totalPages++
	}
	return Result[T]{
		Items:      items,
		Page:       p.Page,
		Limit:      p.Limit,
		TotalCount: totalCount,
		TotalPages: totalPages,
	}
}
