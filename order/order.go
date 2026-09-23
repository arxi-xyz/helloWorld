package order

import "helloworld/pagination"

type PageNumber int
type PageSize int

type Status string

const (
	StatusOpen Status = "open"
	StatusPaid Status = "paid"
)

type Order struct {
	ID         int64
	CustomerID int64
	Status     Status
	TotalCents int64
}

type OrderDTO struct {
	ID         int64  `json:"id"`
	CustomerID int64  `json:"customer_id"`
	Status     string `json:"status"`
	TotalCents int64  `json:"total_cents"`
}

func ToDTO(o Order) OrderDTO {
	return OrderDTO{
		ID:         o.ID,
		CustomerID: o.CustomerID,
		Status:     string(o.Status),
		TotalCents: o.TotalCents,
	}
}

// Repository is not generic. ListByCustomer is an order query with order
// arguments. A Repository[T] would erase that query and push the customer rule
// out of this package. Page[T] is generic only because a page of OrderDTO and a
// page of UserDTO share a shape, not a domain.
type Repository interface {
	ListByCustomer(customerID int64, number PageNumber, size PageSize) (orders []Order, total int)
}

type Memory struct {
	Orders []Order
}

func (m Memory) ListByCustomer(customerID int64, number PageNumber, size PageSize) ([]Order, int) {
	matched := pagination.Filter(m.Orders, func(o Order) bool {
		return o.CustomerID == customerID
	})
	return pagination.Slice(matched, number, size), len(matched)
}

func (m Memory) PageByCustomer(customerID int64, number PageNumber, size PageSize) pagination.Page[OrderDTO] {
	orders, total := m.ListByCustomer(customerID, number, size)
	return pagination.PageOf(pagination.Map(orders, ToDTO), number, size, total)
}
