package billing

import (
	"errors"
	"time"

	"helloworld/subscription"
)

var (
	ErrUnknownPlan   = errors.New("billing: plan has no price")
	ErrInvalidPeriod = errors.New("billing: renewal period is invalid")
)

// Invoice is a billing document built from a subscription.Renewal. It does not
// mutate the subscription.
type Invoice struct {
	plan   string
	amount Money
	from   time.Time
	until  time.Time
}

func ListPrice(plan subscription.Plan) (Money, error) {
	switch plan {
	case subscription.Monthly():
		return USD(1500)
	case subscription.Annual():
		return USD(15000)
	default:
		return Money{}, ErrUnknownPlan
	}
}

func InvoiceForRenewal(renewal subscription.Renewal, amount Money) (Invoice, error) {
	if renewal.Plan().IsZero() || !renewal.From().Before(renewal.Until()) {
		return Invoice{}, ErrInvalidPeriod
	}
	if amount.IsZero() {
		return Invoice{}, ErrInvalidAmount
	}
	return Invoice{
		plan:   renewal.Plan().Name(),
		amount: amount,
		from:   renewal.From(),
		until:  renewal.Until(),
	}, nil
}

func (inv Invoice) Plan() string { return inv.plan }

func (inv Invoice) Amount() Money { return inv.amount }

func (inv Invoice) From() time.Time { return inv.from }

func (inv Invoice) Until() time.Time { return inv.until }
