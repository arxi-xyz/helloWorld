package main

import "testing"

func TestParser(t *testing.T) {
	testCases := []struct {
		name           string
		reqTimeout     string
		maxRetries     string
		minOrderAmount string
		expectedError  error
	}{
		{"valid", "5s", "3", "10.00", nil},
		{"invalid reqTimeout", "invalid", "3", "10.00", nil},
		{"invalid maxRetries", "5s", "invalid", "10.00", nil},
		{"invalid minOrderAmount", "5s", "3", "invalid", nil},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			time, retries, orderAmount, err := parser(tt.reqTimeout, tt.maxRetries, tt.minOrderAmount)

			if err != tt.expectedError {
				t.Errorf("got %d, want %d", err, tt.expectedError)
			}
		})
	}
}
