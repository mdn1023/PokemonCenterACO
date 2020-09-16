package models

// Profile represents the json request body for submitting billing and shipping address info
type Profile struct {
	Billing  Address `json:"billing"`
	Shipping Address `json:"shipping"`
}

// Address represents the json inner structs for billing and shipping addresses
type Address struct {
	FirstName    string `json:"givenName"`
	LastName     string `json:"familyName"`
	Country      string `json:"countryName"`
	City         string `json:"locality"`
	PhoneNumber  string `json:"phoneNumber"`
	PostalCode   string `json:"postalCode"`
	State        string `json:"region"`
	AddressLine1 string `json:"streetAddress"`
	AddressLine2 string `json:"extendedAddress"`
}
