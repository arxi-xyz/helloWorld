package payment_test

import (
	"context"
	"errors"
	"testing"

	"helloworld/payment"
	"helloworld/payment/adapter"
	"helloworld/payment/memory"
	"helloworld/provider"
)

var (
	_ payment.Gateway = (*memory.Gateway)(nil)
	_ payment.Gateway = (*adapter.Adapter)(nil)
)

func TestChargeSuccess(t *testing.T) {
	gateway := &memory.Gateway{}
	svc, err := payment.NewService(gateway)
	if err != nil {
		t.Fatal(err)
	}

	charge := payment.Charge{CustomerID: "cust_1", Amount: 2500, Currency: "USD"}
	receipt, err := svc.Charge(context.Background(), charge)
	if err != nil {
		t.Fatal(err)
	}
	if receipt.ID != "mem_1" {
		t.Fatalf("receipt id = %s", receipt.ID)
	}
	if receipt.Charge != charge {
		t.Fatalf("receipt charge = %+v", receipt.Charge)
	}

	stored := gateway.Charges()
	if len(stored) != 1 || stored[0] != receipt {
		t.Fatalf("stored charges = %+v", stored)
	}
}

func TestChargeProviderError(t *testing.T) {
	client := &provider.Client{
		Err: &provider.Error{Code: "card_declined", Message: "insufficient funds"},
	}
	svc, err := payment.NewService(adapter.New(client))
	if err != nil {
		t.Fatal(err)
	}

	receipt, err := svc.Charge(context.Background(), payment.Charge{
		CustomerID: "cust_1",
		Amount:     2500,
		Currency:   "USD",
	})
	if err == nil {
		t.Fatal("expected provider error")
	}
	var provErr *provider.Error
	if !errors.As(err, &provErr) {
		t.Fatalf("error = %T %v", err, err)
	}
	if provErr.Code != "card_declined" {
		t.Fatalf("code = %s", provErr.Code)
	}
	if receipt != (payment.Receipt{}) {
		t.Fatalf("receipt = %+v, want zero", receipt)
	}
}

func TestTypedNilProviderErrorIsSuccess(t *testing.T) {
	var provErr *provider.Error
	client := &provider.Client{ChargeID: "ch_ok", Err: provErr}
	svc, err := payment.NewService(adapter.New(client))
	if err != nil {
		t.Fatal(err)
	}

	receipt, err := svc.Charge(context.Background(), payment.Charge{
		CustomerID: "cust_1",
		Amount:     2500,
		Currency:   "USD",
	})
	if err != nil {
		t.Fatalf("typed nil *provider.Error surfaced as %T %[1]v", err)
	}
	if receipt.ID != "ch_ok" {
		t.Fatalf("receipt id = %s", receipt.ID)
	}
}
