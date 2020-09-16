package models

type SubmitOrderResponse struct {
	Self struct {
		Type string `json:"type"`
		URI  string `json:"uri"`
		Href string `json:"href"`
	} `json:"self"`
	Messages []struct {
		Type         string `json:"type"`
		ID           string `json:"id"`
		DebugMessage string `json:"debug-message"`
		Data         struct {
			Cause string `json:"cause"`
		} `json:"data"`
	} `json:"messages"`
	Links []struct {
		Rel  string `json:"rel"`
		Type string `json:"type"`
		URI  string `json:"uri"`
		Href string `json:"href"`
		Rev  string `json:"rev,omitempty"`
	} `json:"links"`
	BillingAddress struct {
		Address struct {
			CountryName     string      `json:"country-name"`
			ExtendedAddress interface{} `json:"extended-address"`
			Locality        string      `json:"locality"`
			Organization    interface{} `json:"organization"`
			PhoneNumber     string      `json:"phone-number"`
			PostalCode      string      `json:"postal-code"`
			Region          string      `json:"region"`
			StreetAddress   string      `json:"street-address"`
		} `json:"address"`
		Name struct {
			FamilyName string `json:"family-name"`
			GivenName  string `json:"given-name"`
		} `json:"name"`
	} `json:"billing-address"`
	MonetaryTotal []struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
		Display  string  `json:"display"`
	} `json:"monetary-total"`
	PaymentMeans string `json:"payment-means"`
	PaymentName  string `json:"payment-name"`
	PurchaseDate struct {
		DisplayValue string `json:"display-value"`
		Value        int64  `json:"value"`
	} `json:"purchase-date"`
	PurchaseNumber       string `json:"purchase-number"`
	ShippingDestinations []struct {
		Address struct {
			CountryName     string      `json:"country-name"`
			ExtendedAddress interface{} `json:"extended-address"`
			Locality        string      `json:"locality"`
			Organization    interface{} `json:"organization"`
			PhoneNumber     string      `json:"phone-number"`
			PostalCode      string      `json:"postal-code"`
			Region          string      `json:"region"`
			StreetAddress   string      `json:"street-address"`
		} `json:"address"`
		Name struct {
			FamilyName string `json:"family-name"`
			GivenName  string `json:"given-name"`
		} `json:"name"`
	} `json:"shipping-destinations"`
	ShippingOptions []struct {
		Carrier string `json:"carrier"`
		Cost    []struct {
			Amount   int    `json:"amount"`
			Currency string `json:"currency"`
			Display  string `json:"display"`
		} `json:"cost"`
		DisplayName string `json:"display-name"`
		Name        string `json:"name"`
	} `json:"shipping-options"`
	Status   string `json:"status"`
	TaxTotal struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
		Display  string  `json:"display"`
	} `json:"tax-total"`
	Taxes []struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
		Title    string  `json:"title"`
	} `json:"taxes"`
}
