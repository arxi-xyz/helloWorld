// Package memory is an in-memory Gateway for tests.
package memory

import (
	"context"
	"fmt"
	"sync"

	"helloworld/payment"
)

type Gateway struct {
	mu      sync.Mutex
	nextID  int
	charges []payment.Receipt
}

func (g *Gateway) Charge(ctx context.Context, charge payment.Charge) (payment.Receipt, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.nextID++
	receipt := payment.Receipt{
		ID:     fmt.Sprintf("mem_%d", g.nextID),
		Charge: charge,
	}
	g.charges = append(g.charges, receipt)
	return receipt, nil
}

func (g *Gateway) Charges() []payment.Receipt {
	g.mu.Lock()
	defer g.mu.Unlock()

	out := make([]payment.Receipt, len(g.charges))
	copy(out, g.charges)
	return out
}
