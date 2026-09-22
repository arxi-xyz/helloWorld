package subscription

import "time"

// Clock supplies the current time. A Subscription is composed with a Clock
// instead of calling time.Now directly.
type Clock interface {
	Now() time.Time
}

// SystemClock reads the host clock. It has no mutable state, so Now uses a
// value receiver.
type SystemClock struct{}

func (SystemClock) Now() time.Time { return time.Now() }
