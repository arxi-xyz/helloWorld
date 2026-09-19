package main

import (
	"errors"
	"strconv"
	"time"
)

type Amount int64

func main() {
	parser("5s", "3", "10.00")
}

func parser(reqTimeout, maxRetries, minOrderAmount string) (time.Duration, int, Amount, error) {
	if reqTimeout == "" {
		return 0, 0, 0, errors.New("request timeout is required")
	}

	if maxRetries == "" {
		return 0, 0, 0, errors.New("max retries is required")
	}

	if minOrderAmount == "" {
		return 0, 0, 0, errors.New("minimum order amount is required")
	}

	timeout, err := time.ParseDuration(reqTimeout)
	if err != nil {
		return 0, 0, 0, errors.New("invalid request timeout")
	}

	retries, err := strconv.Atoi(maxRetries)
	if err != nil {
		return 0, 0, 0, errors.New("invalid max retries")
	}

	orderAmountFloat, err := strconv.ParseFloat(minOrderAmount, 64)
	if err != nil {
		return 0, 0, 0, errors.New("invalid minimum order amount")
	}

	orderAmount := Amount(orderAmountFloat * 100)

	return timeout, retries, orderAmount, nil
}
