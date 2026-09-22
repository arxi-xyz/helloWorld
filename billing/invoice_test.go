package billing_test

import (
	"testing"
	"time"

	"helloworld/billing"
	"helloworld/subscription"
)

type fakeClock struct {
	now time.Time
}

func (c *fakeClock) Now() time.Time { return c.now }

func TestInvoiceUsesRenewalFromSubscription(t *testing.T) {
	start := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: start}
	sub, err := subscription.New(subscription.Annual(), clock, start)
	if err != nil {
		t.Fatal(err)
	}

	renewal, err := sub.Renew()
	if err != nil {
		t.Fatal(err)
	}
	price, err := billing.ListPrice(renewal.Plan())
	if err != nil {
		t.Fatal(err)
	}
	invoice, err := billing.InvoiceForRenewal(renewal, price)
	if err != nil {
		t.Fatal(err)
	}

	if invoice.Plan() != "annual" {
		t.Fatalf("plan = %s", invoice.Plan())
	}
	if invoice.Amount().Cents() != 15000 || invoice.Amount().Currency() != "USD" {
		t.Fatalf("amount = %d %s", invoice.Amount().Cents(), invoice.Amount().Currency())
	}
	if !invoice.From().Equal(renewal.From()) || !invoice.Until().Equal(renewal.Until()) {
		t.Fatal("invoice period does not match the renewal")
	}
}

func TestListPriceRejectsUnknownPlan(t *testing.T) {
	if _, err := billing.ListPrice(subscription.Plan{}); err != billing.ErrUnknownPlan {
		t.Fatalf("error = %v", err)
	}
}

func TestInvoiceRejectsInvalidMoneyAndPeriod(t *testing.T) {
	start := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	sub, err := subscription.New(subscription.Monthly(), &fakeClock{now: start}, start)
	if err != nil {
		t.Fatal(err)
	}
	renewal, err := sub.Renew()
	if err != nil {
		t.Fatal(err)
	}

	if _, err := billing.USD(0); err != billing.ErrInvalidAmount {
		t.Fatalf("USD(0) error = %v", err)
	}
	if _, err := billing.InvoiceForRenewal(renewal, billing.Money{}); err != billing.ErrInvalidAmount {
		t.Fatalf("zero money error = %v", err)
	}
	if _, err := billing.InvoiceForRenewal(subscription.Renewal{}, mustUSD(t, 1500)); err != billing.ErrInvalidPeriod {
		t.Fatalf("zero renewal error = %v", err)
	}
}

func mustUSD(t *testing.T, cents int64) billing.Money {
	t.Helper()
	money, err := billing.USD(cents)
	if err != nil {
		t.Fatal(err)
	}
	return money
}
