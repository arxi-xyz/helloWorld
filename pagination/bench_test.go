package pagination_test

import (
	"testing"

	"helloworld/pagination"
)

var benchSink any

type benchRow struct {
	ID   int
	Name string
}

func benchRows(n int) []benchRow {
	rows := make([]benchRow, n)
	for i := range rows {
		rows[i] = benchRow{ID: i, Name: "item"}
	}
	return rows
}

func BenchmarkMapGeneric(b *testing.B) {
	rows := benchRows(256)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchSink = pagination.Map(rows, func(r benchRow) benchRow {
			return benchRow{ID: r.ID, Name: r.Name}
		})
	}
}

func BenchmarkMapConcrete(b *testing.B) {
	rows := benchRows(256)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		out := make([]benchRow, len(rows))
		for j, r := range rows {
			out[j] = benchRow{ID: r.ID, Name: r.Name}
		}
		benchSink = out
	}
}
