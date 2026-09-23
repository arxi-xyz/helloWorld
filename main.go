package main

import (
	"encoding/json"
	"log"
	"os"

	"helloworld/order"
	"helloworld/user"
)

func main() {
	orders := order.Memory{Orders: []order.Order{
		{ID: 1, CustomerID: 7, Status: order.StatusOpen, TotalCents: 1200},
		{ID: 2, CustomerID: 7, Status: order.StatusPaid, TotalCents: 800},
		{ID: 3, CustomerID: 8, Status: order.StatusOpen, TotalCents: 500},
		{ID: 4, CustomerID: 7, Status: order.StatusOpen, TotalCents: 250},
	}}
	users := user.Memory{Users: []user.User{
		{ID: 1, Name: "Ada", Active: true},
		{ID: 2, Name: "Grace", Active: false},
		{ID: 3, Name: "Lin", Active: true},
	}}

	enc := json.NewEncoder(os.Stdout)
	if err := enc.Encode(orders.PageByCustomer(7, 1, 2)); err != nil {
		log.Fatal(err)
	}
	if err := enc.Encode(users.PageActive(1, 10)); err != nil {
		log.Fatal(err)
	}
}
