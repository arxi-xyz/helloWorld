package subscription

import (
	"errors"
	"time"
)

var (
	ErrInvalidPlan      = errors.New("subscription: plan is invalid")
	ErrNoClock          = errors.New("subscription: clock is required")
	ErrInvalidStart     = errors.New("subscription: start time is required")
	ErrCancelled        = errors.New("subscription: cancelled subscription cannot be changed")
	ErrAlreadyCancelled = errors.New("subscription: subscription is already cancelled")
)

// Subscription is a mutable entity, so its methods use pointer receivers.
// The clock field is a composed dependency, not an embedded base type.
type Subscription struct {
	plan      Plan
	expiresAt time.Time
	cancelled bool
	clock     Clock
}

func New(plan Plan, clock Clock, startedAt time.Time) (*Subscription, error) {
	if plan.IsZero() {
		return nil, ErrInvalidPlan
	}
	if clock == nil {
		return nil, ErrNoClock
	}
	if startedAt.IsZero() {
		return nil, ErrInvalidStart
	}
	return &Subscription{
		plan:      plan,
		expiresAt: startedAt.Add(plan.Term()),
		clock:     clock,
	}, nil
}

func (s *Subscription) Plan() Plan { return s.plan }

func (s *Subscription) ExpiresAt() time.Time { return s.expiresAt }

func (s *Subscription) IsCancelled() bool { return s.cancelled }

// IsActive is false when the subscription is cancelled or when the clock is
// at or past the expiration instant.
func (s *Subscription) IsActive() bool {
	if s.cancelled {
		return false
	}
	return s.clock.Now().Before(s.expiresAt)
}

// ChangePlan replaces the plan used by later renewals. It does not move the
// current expiration, and a cancelled subscription stays cancelled.
func (s *Subscription) ChangePlan(plan Plan) error {
	if s.cancelled {
		return ErrCancelled
	}
	if plan.IsZero() {
		return ErrInvalidPlan
	}
	s.plan = plan
	return nil
}

// Cancel is terminal. A second call fails so the cancelled flag is not
// treated as a toggle.
func (s *Subscription) Cancel() error {
	if s.cancelled {
		return ErrAlreadyCancelled
	}
	s.cancelled = true
	return nil
}

// Renew extends a still-active subscription from its current expiration, so
// remaining time is kept. An expired subscription starts the new term at the
// clock's current time. A cancelled subscription cannot renew.
func (s *Subscription) Renew() (Renewal, error) {
	if s.cancelled {
		return Renewal{}, ErrCancelled
	}
	if s.plan.IsZero() {
		return Renewal{}, ErrInvalidPlan
	}

	now := s.clock.Now()
	from := s.expiresAt
	if !now.Before(s.expiresAt) {
		from = now
	}
	until := from.Add(s.plan.Term())
	if !from.Before(until) {
		return Renewal{}, ErrInvalidPlan
	}

	s.expiresAt = until
	return Renewal{plan: s.plan, from: from, until: until}, nil
}

// Renewal is an immutable record of one successful renewal. Billing reads it
// and does not reach into Subscription.
type Renewal struct {
	plan  Plan
	from  time.Time
	until time.Time
}

func (r Renewal) Plan() Plan { return r.plan }

func (r Renewal) From() time.Time { return r.from }

func (r Renewal) Until() time.Time { return r.until }
