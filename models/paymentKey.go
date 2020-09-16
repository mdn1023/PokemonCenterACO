package models

// PaymentKey is a JSON representation of the payment key received in the response body after successfully submitting billing info
type PaymentKey struct {
	KeyID string `json:"keyId"`
}
