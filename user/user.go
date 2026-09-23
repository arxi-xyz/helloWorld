package user

import "helloworld/pagination"

type PageNumber int
type PageSize int

type User struct {
	ID     int64
	Name   string
	Active bool
}

type UserDTO struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

func ToDTO(u User) UserDTO {
	return UserDTO{ID: u.ID, Name: u.Name, Active: u.Active}
}

// Repository stays concrete for the same reason as order.Repository: ListActive
// is a user rule, not a generic list of T.
type Repository interface {
	ListActive(number PageNumber, size PageSize) (users []User, total int)
}

type Memory struct {
	Users []User
}

func (m Memory) ListActive(number PageNumber, size PageSize) ([]User, int) {
	matched := pagination.Filter(m.Users, func(u User) bool { return u.Active })
	return pagination.Slice(matched, number, size), len(matched)
}

func (m Memory) PageActive(number PageNumber, size PageSize) pagination.Page[UserDTO] {
	users, total := m.ListActive(number, size)
	return pagination.PageOf(pagination.Map(users, ToDTO), number, size, total)
}
