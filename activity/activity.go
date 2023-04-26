package app

import (
	"context"
	"errors"
	"fmt"
	mo "temporal-tutorial/model"
	"time"
)

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

func MyActivity2(ctx context.Context, number int) (string, error) {
	if number > 12 {
		return "", &mo.BadRequestError{}
	}
	m := time.Month(number)
	greeting := fmt.Sprintf("Sekarang Bulan %s!", m)
	return greeting, nil
}

func MyActivity3(ctx context.Context, number1, number2 int) (string, error) {
	greeting := fmt.Sprintf("Perkalian %d x %d = %d", number1, number2, (number1 * number2))
	return greeting, nil
}
