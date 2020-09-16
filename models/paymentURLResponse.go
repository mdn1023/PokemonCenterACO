package models

// PaymentURLResponse is a struct that represents the response body json received when requesting the payment url
type PaymentURLResponse struct {
	Self struct {
		Type string `json:"type"`
		URI  string `json:"uri"`
		Href string `json:"href"`
	} `json:"self"`
	Messages    []interface{} `json:"messages"`
	Links       []interface{} `json:"links"`
	DisplayName string        `json:"display-name"`
	Token       string        `json:"token"`
}
