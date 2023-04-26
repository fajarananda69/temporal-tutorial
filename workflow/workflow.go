package app

import (
	"net/http"
	"strconv"
	ac "temporal-tutorial/activity"
	mo "temporal-tutorial/model"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

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

func MyWorkflow1(ctx workflow.Context, service string, input string) (mo.Response, error) {
	var err error
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retryPolice(),
		ActivityID:  "Test Activity",
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	if service == "1" {
		err = workflow.ExecuteActivity(ctx, ac.MyActivity1, input).Get(ctx, &result)
	} else {
		number, _ := strconv.Atoi(input)
		err = workflow.ExecuteActivity(ctx, ac.MyActivity2, number).Get(ctx, &result)
	}

	response := mo.Response{
		Status: http.StatusOK,
		Data:   result,
	}
	return response, err
}

func MyWorkflow2(ctx workflow.Context, service string, input1 string, input2 string) (mo.Response, error) {
	var err error
	options := workflow.ActivityOptions{
		// Timeout options specify when to automatically timeout Activity functions.
		StartToCloseTimeout: time.Minute,
		// Optionally provide a customized RetryPolicy.
		// Temporal retries failed Activities by default.
		RetryPolicy: retryPolice(),
		ActivityID:  "Test Activity",
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	var result string
	number1, _ := strconv.Atoi(input1)
	number2, _ := strconv.Atoi(input2)
	err = workflow.ExecuteActivity(ctx, ac.MyActivity3, number1, number2).Get(ctx, &result)

	response := mo.Response{
		Status: http.StatusOK,
		Data:   result,
	}
	return response, err
}
