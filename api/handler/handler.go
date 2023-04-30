package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	mo "temporal-tutorial/model"
	wo "temporal-tutorial/workflow"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"go.temporal.io/sdk/client"
)

var (
	Temporal client.Client
)

func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	res := make(map[string]interface{})
	res["products"] = mo.Products

	Yay(w, r, http.StatusOK, res)

}

func CreateCartHandler(w http.ResponseWriter, r *http.Request) {
	workflowID := "MY_CART-" + fmt.Sprintf("%d", time.Now().Unix())

	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: mo.MyTaskQueue2,
	}

	cart := mo.CartState{Items: make([]mo.CartItem, 0)}
	we, err := Temporal.ExecuteWorkflow(context.Background(), options, wo.CartWorkflow, cart)
	if err != nil {
		Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	res := make(map[string]interface{})
	res["workflowID"] = we.GetID()

	Yay(w, r, http.StatusCreated, res)
}

func GetCartHandler(w http.ResponseWriter, r *http.Request) {
	var (
		workflowID = chi.URLParam(r, "workflowID")
		ctx        = r.Context()
	)

	response, err := Temporal.QueryWorkflow(ctx, workflowID, "", mo.MyQuery)
	if err != nil {
		Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	var res interface{}
	if err := response.Get(&res); err != nil {
		Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	Yay(w, r, http.StatusOK, res)
}

func AddToCartHandler(w http.ResponseWriter, r *http.Request) {
	var (
		workflowID = chi.URLParam(r, "workflowID")
		ctx        = r.Context()
	)

	var item mo.CartItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	update := mo.AddToCartSignal{Route: mo.ROUTE_ADD_TO_CART, Item: item}

	err = Temporal.SignalWorkflow(ctx, workflowID, "", mo.SIGNAL_ADD_TO_CART_CHANNEL, update)
	if err != nil {
		Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	res := map[string]interface{}{
		"ok":       1,
		"add_item": item,
	}

	Yay(w, r, http.StatusOK, res)
}

func RemoveFromCartHandler(w http.ResponseWriter, r *http.Request) {
	var (
		workflowID = chi.URLParam(r, "workflowID")
		ctx        = r.Context()
	)

	var item mo.CartItem
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	update := mo.AddToCartSignal{Route: mo.ROUTE_REMOVE_FROM_CART, Item: item}

	err = Temporal.SignalWorkflow(ctx, workflowID, "", mo.SIGNAL_REMOVE_FROM_CART_CHANNEL, update)
	if err != nil {
		Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	res := map[string]interface{}{
		"ok":          1,
		"remove_item": item,
	}

	Yay(w, r, http.StatusOK, res)
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) {
	var (
		workflowID = chi.URLParam(r, "workflowID")
		ctx        = r.Context()
	)

	checkout := mo.CheckoutSignal{Route: mo.ROUTE_CHECKOUT}

	err := Temporal.SignalWorkflow(ctx, workflowID, "", mo.SIGNAL_CHECKOUT_CHANNEL, checkout)
	if err != nil {
		Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	res := map[string]interface{}{
		"checkout": true,
	}

	Yay(w, r, http.StatusOK, res)
}

func PaymentHandler(w http.ResponseWriter, r *http.Request) {
	var (
		workflowID = chi.URLParam(r, "workflowID")
		ctx        = r.Context()
	)

	payment := mo.CheckoutSignal{Route: mo.ROUTE_PAYMENT}

	err := Temporal.SignalWorkflow(ctx, workflowID, "", mo.SIGNAL_PAYMENT_CHANNEL, payment)
	if err != nil {
		Nay(w, r, http.StatusInternalServerError, err)
		return
	}

	res := map[string]interface{}{
		"payment": true,
	}

	Yay(w, r, http.StatusOK, res)
}

func Yay(w http.ResponseWriter, r *http.Request, status int, content interface{}) {
	render.Status(r, status)
	_ = render.Render(w, r, &mo.Response{
		Data:   content,
		Status: status,
	})
}

func Nay(w http.ResponseWriter, r *http.Request, status int, err error) {
	render.Status(r, status)
	_ = render.Render(w, r, &mo.Response{
		Status: status,
		Error:  err,
	})
}
