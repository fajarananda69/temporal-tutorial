package app

import (
	"context"
	mo "temporal-tutorial/model"

	"github.com/segmentio/ksuid"
)

type Activities struct{}

// checkout activity
func (a *Activities) Checkout(ctx context.Context, state *mo.CartState) (mo.Transaction, error) {
	if len(state.Items) == 0 {
		return mo.Transaction{}, &mo.CartEmptyError{}
	}
	status := mo.Transaction{
		TrxID: ksuid.New().String(),
	}
	return status, nil
}

// payment activity
func (a *Activities) Payment(ctx context.Context, state *mo.CartState) (mo.Transaction, error) {
	if !state.Checkout {
		return mo.Transaction{}, &mo.CartEmptyError{}
	}

	status := mo.Transaction{
		PaymentID: ksuid.New().String(),
	}
	return status, nil
}
