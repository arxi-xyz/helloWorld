package payment

import (
	"context"
	"errors"
)

var ErrNoGateway = errors.New("payment: gateway is required")

type Charge struct {
	CustomerID string
	Amount     int64
	Currency   string
}

type Receipt struct {
	ID string
	Charge
}

// Gateway is the only provider capability Service needs. Implementations live
// outside this package.
type Gateway interface {
	Charge(ctx context.Context, charge Charge) (Receipt, error)
}

// Service charges through a composed Gateway. Pointer receivers keep the
// gateway reference with the service instead of copying it.
type Service struct {
	gateway Gateway
}

func NewService(gateway Gateway) (*Service, error) {
	if gateway == nil {
		return nil, ErrNoGateway
	}
	return &Service{gateway: gateway}, nil
}

func (s *Service) Charge(ctx context.Context, charge Charge) (Receipt, error) {
	return s.gateway.Charge(ctx, charge)
}
