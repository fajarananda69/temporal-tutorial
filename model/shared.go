package model

import "net/http"

// set Error
type BadRequestError struct{}
type CartEmptyError struct{}

func (m *BadRequestError) Error() string {
	return "Request is invalid"
}

func (m *CartEmptyError) Error() string {
	return "Cart is empty"
}

// Task Queue
const (
	MyTaskQueue1 = "MY_TASK_QUEUE_1"
	MyTaskQueue2 = "MY_TASK_QUEUE_2"
)

const (
	MyQuery = "getCart"
)

// Signal
const (
	SIGNAL_ADD_TO_CART_CHANNEL      = "SIGNAL_ADD_TO_CART_CHANNEL"
	SIGNAL_REMOVE_FROM_CART_CHANNEL = "SIGNAL_REMOVE_FROM_CART_CHANNEL"
	SIGNAL_CHECKOUT_CHANNEL         = "SIGNAL_CHECKOUT_CHANNEL"
	SIGNAL_PAYMENT_CHANNEL          = "SIGNAL_PAYMENT_CHANNEL"
)

const (
	ROUTE_ADD_TO_CART      = "ROUTE_ADD_TO_CART"
	ROUTE_REMOVE_FROM_CART = "ROUTE_REMOVE_FROM_CART"
	ROUTE_CHECKOUT         = "ROUTE_CHECKOUT"
	ROUTE_PAYMENT          = "ROUTE_PAYMENT"
)

// response
type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  error       `json:"error,omitempty"`
}

func (res *Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
