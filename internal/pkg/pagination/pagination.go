package pagination

const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

type Params struct {
	Page     int
	PageSize int
}

type Meta struct {
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
}

func Normalize(page int, pageSize int) Params {
	if page < 1 {
		page = DefaultPage
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	return Params{Page: page, PageSize: pageSize}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func NewMeta(params Params, totalItems int) Meta {
	totalPages := 0
	if totalItems > 0 {
		totalPages = (totalItems + params.PageSize - 1) / params.PageSize
	}
	return Meta{
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}
