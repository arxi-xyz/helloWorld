package subscription_test

import (
	"testing"
	"time"

	"helloworld/subscription"
)

// fakeClock is a hand-written Clock. Advance moves time without a mock library.
type fakeClock struct {
	now time.Time
}

func (c *fakeClock) Now() time.Time { return c.now }

func (c *fakeClock) Advance(d time.Duration) { c.now = c.now.Add(d) }

func TestNewRejectsBrokenInputs(t *testing.T) {
	clock := &fakeClock{now: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	start := clock.Now()

	if _, err := subscription.New(subscription.Plan{}, clock, start); err != subscription.ErrInvalidPlan {
		t.Fatalf("zero plan error = %v", err)
	}
	if _, err := subscription.New(subscription.Monthly(), nil, start); err != subscription.ErrNoClock {
		t.Fatalf("nil clock error = %v", err)
	}
	if _, err := subscription.New(subscription.Monthly(), clock, time.Time{}); err != subscription.ErrInvalidStart {
		t.Fatalf("zero start error = %v", err)
	}
}

func TestExpirationFollowsFakeClock(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: start}
	sub, err := subscription.New(subscription.Monthly(), clock, start)
	if err != nil {
		t.Fatal(err)
	}

	if !sub.IsActive() {
		t.Fatal("subscription should be active at the start")
	}
	wantExpiry := start.Add(subscription.Monthly().Term())
	if !sub.ExpiresAt().Equal(wantExpiry) {
		t.Fatalf("expires at %s, want %s", sub.ExpiresAt(), wantExpiry)
	}

	clock.Advance(subscription.Monthly().Term() - time.Second)
	if !sub.IsActive() {
		t.Fatal("subscription should stay active one second before expiration")
	}

	clock.Advance(time.Second)
	if sub.IsActive() {
		t.Fatal("subscription should expire when the fake clock reaches ExpiresAt")
	}
}

func TestRenewKeepsRemainingTimeThenRestartsAfterExpiry(t *testing.T) {
	start := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: start}
	sub, err := subscription.New(subscription.Monthly(), clock, start)
	if err != nil {
		t.Fatal(err)
	}
	originalExpiry := sub.ExpiresAt()

	clock.Advance(10 * 24 * time.Hour)
	renewal, err := sub.Renew()
	if err != nil {
		t.Fatal(err)
	}
	if !renewal.From().Equal(originalExpiry) {
		t.Fatalf("active renewal starts at %s, want current expiry %s", renewal.From(), originalExpiry)
	}
	if !renewal.Until().Equal(originalExpiry.Add(subscription.Monthly().Term())) {
		t.Fatalf("active renewal ends at %s", renewal.Until())
	}
	if !sub.ExpiresAt().Equal(renewal.Until()) {
		t.Fatal("renewal did not move expiration")
	}

	clock.now = sub.ExpiresAt()
	if sub.IsActive() {
		t.Fatal("subscription should be expired at the new expiration")
	}
	expiredRenewal, err := sub.Renew()
	if err != nil {
		t.Fatal(err)
	}
	if !expiredRenewal.From().Equal(clock.Now()) {
		t.Fatalf("expired renewal starts at %s, want clock %s", expiredRenewal.From(), clock.Now())
	}
	if !expiredRenewal.Until().Equal(clock.Now().Add(subscription.Monthly().Term())) {
		t.Fatalf("expired renewal ends at %s", expiredRenewal.Until())
	}
}

func TestChangePlanAppliesOnNextRenewal(t *testing.T) {
	start := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: start}
	sub, err := subscription.New(subscription.Monthly(), clock, start)
	if err != nil {
		t.Fatal(err)
	}
	expiryBeforeChange := sub.ExpiresAt()

	if err := sub.ChangePlan(subscription.Annual()); err != nil {
		t.Fatal(err)
	}
	if !sub.ExpiresAt().Equal(expiryBeforeChange) {
		t.Fatal("changing plan moved the current expiration")
	}

	renewal, err := sub.Renew()
	if err != nil {
		t.Fatal(err)
	}
	if renewal.Plan() != subscription.Annual() {
		t.Fatalf("renewed plan = %s", renewal.Plan().Name())
	}
	if !renewal.Until().Equal(expiryBeforeChange.Add(subscription.Annual().Term())) {
		t.Fatalf("renewed until %s", renewal.Until())
	}
}

func TestCancelBlocksRenewalAndPlanChange(t *testing.T) {
	start := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: start}
	sub, err := subscription.New(subscription.Monthly(), clock, start)
	if err != nil {
		t.Fatal(err)
	}

	if err := sub.Cancel(); err != nil {
		t.Fatal(err)
	}
	if sub.IsActive() {
		t.Fatal("cancelled subscription should not be active")
	}
	expiry := sub.ExpiresAt()

	if _, err := sub.Renew(); err != subscription.ErrCancelled {
		t.Fatalf("renew after cancel error = %v", err)
	}
	if err := sub.ChangePlan(subscription.Annual()); err != subscription.ErrCancelled {
		t.Fatalf("change plan after cancel error = %v", err)
	}
	if err := sub.Cancel(); err != subscription.ErrAlreadyCancelled {
		t.Fatalf("second cancel error = %v", err)
	}
	if !sub.ExpiresAt().Equal(expiry) || sub.Plan() != subscription.Monthly() {
		t.Fatal("failed changes mutated the cancelled subscription")
	}
}
