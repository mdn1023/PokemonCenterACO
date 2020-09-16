package models

type PaymentInfo struct {
	Name             string
	CreditCardType   string
	CreditCardNumber string
	Cvv              string
	ExpiryMonth      string
	ExpiryYear       string
}
