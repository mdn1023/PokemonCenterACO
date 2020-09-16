package models

// ATC is a json representation of the response body returned from the POST request for adding and item to cart
type ATC struct {
	Messages []struct {
		Type         string `json:"type"`
		ID           string `json:"id"`
		DebugMessage string `json:"debug-message"`
		Data         struct {
			ItemCode string `json:"item-code"`
		} `json:"data"`
	} `json:"messages"`
	Links []interface{} `json:"links"`
}
