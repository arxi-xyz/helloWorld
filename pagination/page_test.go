package pagination_test

import (
	"testing"

	"helloworld/pagination"
)

func TestPageOfAcceptsNamedIntegers(t *testing.T) {
	type PageNumber int
	type PageSize int64

	page := pagination.PageOf([]string{"a", "b", "c"}, PageNumber(0), PageSize(2), 3)
	if page.Number != 1 || page.Size != 2 || page.Total != 3 || page.TotalPages != 2 {
		t.Fatalf("page = %+v", page)
	}
	if len(page.Items) != 3 || page.Items[0] != "a" {
		t.Fatalf("items = %#v", page.Items)
	}
}

func TestSliceUsesNamedIndexes(t *testing.T) {
	type PageNumber int
	type RowLimit int64

	got := pagination.Slice([]int{1, 2, 3, 4, 5}, PageNumber(2), RowLimit(2))
	if len(got) != 2 || got[0] != 3 || got[1] != 4 {
		t.Fatalf("slice = %v", got)
	}

	if tail := pagination.Slice([]int{1}, PageNumber(2), RowLimit(10)); len(tail) != 0 {
		t.Fatalf("out of range slice = %v", tail)
	}
}

func TestMapAndFilter(t *testing.T) {
	type cents int64

	got := pagination.Map([]cents{10, 20}, func(v cents) int64 { return int64(v) * 2 })
	if len(got) != 2 || got[0] != 20 || got[1] != 40 {
		t.Fatalf("map = %v", got)
	}

	kept := pagination.Filter([]cents{10, 15, 20}, func(v cents) bool { return v >= 15 })
	if len(kept) != 2 || kept[0] != 15 || kept[1] != 20 {
		t.Fatalf("filter = %v", kept)
	}
}
