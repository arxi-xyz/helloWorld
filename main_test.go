package main

import (
	"testing"
)

func TestBatchProcess_NilInput(t *testing.T) {
	var orders []Order

	result := batchProcess(orders)

	if result == nil {
		t.Fatal("expected non-nil batch result")
	}

	if len(result) != 0 {
		t.Fatalf("expected 0 batches, got %d", len(result))
	}
}

func TestBatchProcess_EmptyInput(t *testing.T) {
	orders := []Order{}

	result := batchProcess(orders)

	if result == nil {
		t.Fatal("expected non-nil batch result")
	}

	if len(result) != 0 {
		t.Fatalf("expected 0 batches, got %d", len(result))
	}
}

func TestBatchProcess_Aliasing(t *testing.T) {
	input := []Order{
		"order1",
		"order2",
		"order3",
	}

	batches := batchProcess(input)

	// mutate original input
	input[0] = "changed"

	if batches[0][0] != "order1" {
		t.Fatalf(
			"aliasing detected: batch changed after input mutation, got %s",
			batches[0][0],
		)
	}
}

func TestBatchProcess_BatchSize(t *testing.T) {
	orders := make([]Order, 250)

	for i := range orders {
		orders[i] = Order("order")
	}

	batches := batchProcess(orders)

	if len(batches) != 3 {
		t.Fatalf("expected 3 batches, got %d", len(batches))
	}

	if len(batches[0]) != 100 {
		t.Fatalf("expected first batch size 100")
	}

	if len(batches[1]) != 100 {
		t.Fatalf("expected second batch size 100")
	}

	if len(batches[2]) != 50 {
		t.Fatalf("expected third batch size 50")
	}
}

func BenchmarkBatchProcess_CopyBoundary(b *testing.B) {
	orders := make([]Order, 10000)

	for i := range orders {
		orders[i] = Order("order")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = batchProcess(orders)
	}
}
