package pagination

// integer accepts named types whose underlying type is int or int64.
type integer interface {
	~int | ~int64
}

const DefaultSize = 20

// Page is one window of T. Items holds that window, not the whole collection.
type Page[T any] struct {
	Items      []T `json:"items"`
	Number     int `json:"number"`
	Size       int `json:"size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// PageOf records paging metadata. number and size may be named integer types.
func PageOf[T any, N integer, S integer](items []T, number N, size S, total int) Page[T] {
	n, s := bounds(number, size)
	if total < 0 {
		total = 0
	}
	copied := make([]T, len(items))
	copy(copied, items)
	return Page[T]{
		Items:      copied,
		Number:     n,
		Size:       s,
		Total:      total,
		TotalPages: totalPages(total, s),
	}
}

// Slice returns the requested window. An out-of-range number yields an empty slice.
func Slice[T any, N integer, S integer](items []T, number N, size S) []T {
	n, s := bounds(number, size)
	start := (n - 1) * s
	if start >= len(items) || start < 0 {
		return []T{}
	}
	end := start + s
	if end > len(items) {
		end = len(items)
	}
	out := make([]T, end-start)
	copy(out, items[start:end])
	return out
}

func Map[T, U any](items []T, fn func(T) U) []U {
	out := make([]U, len(items))
	for i, item := range items {
		out[i] = fn(item)
	}
	return out
}

func Filter[T any](items []T, keep func(T) bool) []T {
	out := make([]T, 0, len(items))
	for _, item := range items {
		if keep(item) {
			out = append(out, item)
		}
	}
	return out
}

func bounds[N integer, S integer](number N, size S) (int, int) {
	n := int(number)
	s := int(size)
	if n < 1 {
		n = 1
	}
	if s < 1 {
		s = DefaultSize
	}
	return n, s
}

func totalPages(total, size int) int {
	if total <= 0 || size <= 0 {
		return 0
	}
	return (total + size - 1) / size
}
