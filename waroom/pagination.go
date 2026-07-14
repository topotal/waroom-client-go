package waroom

import (
	"iter"

	gen "github.com/topotal/waroom-client-go/internal/ogen"
)

// Pagination holds offset-pagination metadata common to all list
// responses. It is embedded in each generated XxxList type, so metadata
// is accessible directly, e.g. list.TotalCount or list.NextPage.
type Pagination struct {
	CurrentPage int
	NextPage    *int // nil on the last page
	PrevPage    *int // nil on the first page
	TotalPages  int
	TotalCount  int
}

// responseMetadata is the shape every list operation's response_metadata
// satisfies through the getters ogen generates.
type responseMetadata interface {
	GetCurrentPage() int
	GetNextPage() gen.NilInt
	GetPrevPage() gen.NilInt
	GetTotalPages() int
	GetTotalCount() int
}

// newPagination converts ogen response metadata; called by generated code.
func newPagination(m responseMetadata) *Pagination {
	return &Pagination{
		CurrentPage: m.GetCurrentPage(),
		NextPage:    nilIntPtr(m.GetNextPage()),
		PrevPage:    nilIntPtr(m.GetPrevPage()),
		TotalPages:  m.GetTotalPages(),
		TotalCount:  m.GetTotalCount(),
	}
}

func nilIntPtr(n gen.NilInt) *int {
	if v, ok := n.Get(); ok {
		return &v
	}
	return nil
}

// paginate drives a page-fetching function from the start page through the
// last page; called by generated ListXxxAll methods.
func paginate[T any](start gen.OptInt, fetch func(page int) ([]T, *Pagination, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		page := start.Or(1)
		for {
			items, p, err := fetch(page)
			if err != nil {
				var zero T
				yield(zero, err)
				return
			}
			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}
			if p == nil || p.NextPage == nil {
				return
			}
			page = *p.NextPage
		}
	}
}
