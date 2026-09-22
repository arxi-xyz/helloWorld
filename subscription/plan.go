package subscription

import "time"

// Plan is an immutable value: name and term are fixed at construction, and
// every method uses a value receiver.
type Plan struct {
	name string
	term time.Duration
}

func Monthly() Plan {
	return Plan{name: "monthly", term: 30 * 24 * time.Hour}
}

func Annual() Plan {
	return Plan{name: "annual", term: 365 * 24 * time.Hour}
}

func (p Plan) Name() string { return p.name }

func (p Plan) Term() time.Duration { return p.term }

func (p Plan) IsZero() bool {
	return p.name == "" || p.term <= 0
}
