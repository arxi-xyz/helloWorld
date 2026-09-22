package billing

import "errors"

var ErrInvalidAmount = errors.New("billing: amount must be positive")

// Money is an immutable value, so its methods use value receivers.
type Money struct {
	cents    int64
	currency string
}

func USD(cents int64) (Money, error) {
	if cents <= 0 {
		return Money{}, ErrInvalidAmount
	}
	return Money{cents: cents, currency: "USD"}, nil
}

func (m Money) Cents() int64 { return m.cents }

func (m Money) Currency() string { return m.currency }

func (m Money) IsZero() bool {
	return m.cents <= 0 || m.currency == ""
}
