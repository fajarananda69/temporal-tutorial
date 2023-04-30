package app

import (
	"context"
	"errors"
	"fmt"
	mo "temporal-tutorial/model"
	"time"
)

// test activity 1
func MyActivity1(ctx context.Context, name string) (string, error) {
	if name == "err" {
		return "", errors.New("testing error")
	}

	if name == "bad" {
		return "", &mo.BadRequestError{}
	}

	greeting := fmt.Sprintf("Hello %s!", name)
	return greeting, nil
}

// test activity 2
func MyActivity2(ctx context.Context, number int) (string, error) {
	if number > 12 {
		return "", &mo.BadRequestError{}
	}
	m := time.Month(number)
	greeting := fmt.Sprintf("Sekarang Bulan %s!", m)
	return greeting, nil
}
