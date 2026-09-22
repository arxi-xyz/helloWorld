// Package provider is a stand-in for an external payment SDK.
// CreateCharge returns *Error, which is nil on success.
package provider

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Code + ": " + e.Message
}

type Client struct {
	ChargeID string
	Err      *Error
}

func (c *Client) CreateCharge(amount int64, currency, customerID string) (string, *Error) {
	if c.Err != nil {
		return "", c.Err
	}
	if c.ChargeID == "" {
		return "ch_provider", nil
	}
	return c.ChargeID, nil
}
