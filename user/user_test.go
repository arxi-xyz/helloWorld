package user_test

import (
	"testing"

	"helloworld/user"
)

func TestPageActive(t *testing.T) {
	repo := user.Memory{Users: []user.User{
		{ID: 1, Name: "Ada", Active: true},
		{ID: 2, Name: "Grace", Active: false},
		{ID: 3, Name: "Lin", Active: true},
	}}

	page := repo.PageActive(1, 10)
	if page.Total != 2 || page.TotalPages != 1 || len(page.Items) != 2 {
		t.Fatalf("page = %+v", page)
	}
	if page.Items[0] != (user.UserDTO{ID: 1, Name: "Ada", Active: true}) {
		t.Fatalf("first = %+v", page.Items[0])
	}
	if page.Items[1].Name != "Lin" {
		t.Fatalf("second = %+v", page.Items[1])
	}

	second := repo.PageActive(2, 1)
	if len(second.Items) != 1 || second.Items[0].ID != 3 || second.Total != 2 {
		t.Fatalf("second page = %+v", second)
	}
}
