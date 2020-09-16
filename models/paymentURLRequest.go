package models

// PaymentURLRequest represents the json struct of the request body required to fetch the final payment url
type PaymentURLRequest struct {
	PaymentDisplay string `json:"paymentDisplay"`
	PaymentKey     string `json:"paymentKey"`
	PaymentToken   string `json:"paymentToken"`
}
