// Package adapter turns the external provider client into a payment.Gateway.
package adapter

import (
	"context"

	"helloworld/payment"
	"helloworld/provider"
)

type Adapter struct {
	client *provider.Client
}

func New(client *provider.Client) *Adapter {
	return &Adapter{client: client}
}

func (a *Adapter) Charge(ctx context.Context, charge payment.Charge) (payment.Receipt, error) {
	id, provErr := a.client.CreateCharge(charge.Amount, charge.Currency, charge.CustomerID)
	if provErr != nil {
		return payment.Receipt{}, provErr
	}
	// A nil *provider.Error must not be returned as error. That conversion
	// boxes a typed nil and makes the caller see a non-nil error on success.
	return payment.Receipt{ID: id, Charge: charge}, nil
}
