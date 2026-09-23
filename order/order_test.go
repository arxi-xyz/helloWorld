package order_test

import (
	"testing"

	"helloworld/order"
)

func TestPageByCustomer(t *testing.T) {
	repo := order.Memory{Orders: []order.Order{
		{ID: 1, CustomerID: 7, Status: order.StatusOpen, TotalCents: 1200},
		{ID: 2, CustomerID: 7, Status: order.StatusPaid, TotalCents: 800},
		{ID: 3, CustomerID: 8, Status: order.StatusOpen, TotalCents: 500},
		{ID: 4, CustomerID: 7, Status: order.StatusOpen, TotalCents: 250},
	}}

	page := repo.PageByCustomer(7, 1, 2)
	if page.Total != 3 || page.TotalPages != 2 || len(page.Items) != 2 {
		t.Fatalf("page = %+v", page)
	}
	if page.Items[0].ID != 1 || page.Items[1].ID != 2 {
		t.Fatalf("items = %+v", page.Items)
	}
	if page.Items[0].Status != "open" || page.Items[1].TotalCents != 800 {
		t.Fatalf("dto = %+v", page.Items)
	}

	second := repo.PageByCustomer(7, 2, 2)
	if len(second.Items) != 1 || second.Items[0].ID != 4 {
		t.Fatalf("second page = %+v", second.Items)
	}
}
