package app

import (
	"net/http"
	"strconv"
	ac "temporal-tutorial/activity"
	mo "temporal-tutorial/model"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type Workflows struct{}

// Define Retry Policies
func retryPolice() *temporal.RetryPolicy {
	// RetryPolicy specifies how to automatically handle retries if an Activity fails.
	return &temporal.RetryPolicy{
		InitialInterval:        time.Second * 5,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        5, // 5 max retries
		NonRetryableErrorTypes: []string{"BadRequestError"},
	}
}

// Define Workflow
func MyWorkflow1(ctx workflow.Context, input string) (mo.Response, error) {
	var err error

	// set retry policies
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retryPolice(),
		ActivityID:  "Test Client Activity",
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var number int
	if n, err := strconv.Atoi(input); err == nil {
		number = n
	}

	// executed activity
	var result string
	if number > 0 {
		err = workflow.ExecuteActivity(ctx, ac.MyActivity2, number).Get(ctx, &result)
	} else {
		err = workflow.ExecuteActivity(ctx, ac.MyActivity1, input).Get(ctx, &result)
	}

	response := mo.Response{
		Status: http.StatusOK,
		Data:   result,
	}
	return response, err
}

// Define Workflow With Signal
func CartWorkflow(ctx workflow.Context, state mo.CartState) (err error) {
	// set option retry policies
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retryPolice(),
		ActivityID:  "Test Cart Activity",
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// set logger
	logger := workflow.GetLogger(ctx)

	// define interface
	var a *ac.Activities
	var w Workflows

	// set Query
	err = workflow.SetQueryHandler(ctx, mo.MyQuery, func(input []byte) (mo.CartState, error) {
		return state, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed.", "Error", err)
		return err
	}

	// handle Signal
	addToCartChannel := workflow.GetSignalChannel(ctx, mo.SIGNAL_ADD_TO_CART_CHANNEL)
	removeFromCartChannel := workflow.GetSignalChannel(ctx, mo.SIGNAL_REMOVE_FROM_CART_CHANNEL)
	checkoutChannel := workflow.GetSignalChannel(ctx, mo.SIGNAL_CHECKOUT_CHANNEL)
	paymentChannel := workflow.GetSignalChannel(ctx, mo.SIGNAL_PAYMENT_CHANNEL)

	for {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(addToCartChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			if state.Checkout {
				logger.Error("Already checkout cannot add/remove items")
				return
			}

			var message mo.AddToCartSignal
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			w.AddToCart(&state, message.Item)
		})

		selector.AddReceive(removeFromCartChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			if state.Checkout {
				logger.Error("Already checkout cannot add/remove items")
				return
			}

			var message mo.RemoveFromCartSignal
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			w.RemoveFromCart(&state, message.Item)
		})

		selector.AddReceive(checkoutChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var message mo.CheckoutSignal
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			ao := workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute,
			}

			ctx = workflow.WithActivityOptions(ctx, ao)

			var trx mo.Transaction
			err = workflow.ExecuteActivity(ctx, a.Checkout, &state).Get(ctx, &trx)
			if err != nil {
				logger.Error("Error checkout process: %v", err)
				return
			}

			state.TrxID = trx.TrxID
			state.Checkout = true
		})

		selector.AddReceive(paymentChannel, func(c workflow.ReceiveChannel, _ bool) {
			var signal interface{}
			c.Receive(ctx, &signal)

			var message mo.PaymentSignal
			err := mapstructure.Decode(signal, &message)
			if err != nil {
				logger.Error("Invalid signal type %v", err)
				return
			}

			ao := workflow.ActivityOptions{
				StartToCloseTimeout: time.Minute,
			}

			ctx = workflow.WithActivityOptions(ctx, ao)

			var trx mo.Transaction
			err = workflow.ExecuteActivity(ctx, a.Payment, &state).Get(ctx, &trx)
			if err != nil {
				logger.Error("Error payment process: %v", err)
				return
			}

			state.PaymentID = trx.PaymentID
			state.Payment = true
		})

		selector.Select(ctx)

		// if checkout = true && payment == true
		// complate the workflow
		if state.Checkout && state.Payment {
			break
		}
	}

	return nil
}

func (w *Workflows) AddToCart(state *mo.CartState, item mo.CartItem) {
	for i := range state.Items {
		if state.Items[i].ProductId != item.ProductId {
			continue
		}

		state.Items[i].Quantity += item.Quantity
		return
	}

	state.Items = append(state.Items, item)
}

func (w *Workflows) RemoveFromCart(state *mo.CartState, item mo.CartItem) {
	for i := range state.Items {
		if state.Items[i].ProductId != item.ProductId {
			continue
		}

		state.Items[i].Quantity -= item.Quantity
		if state.Items[i].Quantity <= 0 {
			state.Items = append(state.Items[:i], state.Items[i+1:]...)
		}
		break
	}
}
