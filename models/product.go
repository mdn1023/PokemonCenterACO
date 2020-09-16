package models

// Product is a json representation of the response body returned from querying a product by product SKU
type Product struct {
	Availability []struct {
		Self struct {
			Type string `json:"type"`
			URI  string `json:"uri"`
			Href string `json:"href"`
		} `json:"self"`
		Messages []interface{} `json:"messages"`
		Links    []interface{} `json:"links"`
		State    string        `json:"state"`
	} `json:"_availability"`
	Items []struct {
		Element []struct {
			AddToCartForm []struct {
				Self struct {
					Type string `json:"type"`
					URI  string `json:"uri"`
					Href string `json:"href"`
				} `json:"self"`
				Messages []interface{} `json:"messages"`
				Links    []struct {
					Rel  string `json:"rel"`
					Type string `json:"type"`
					URI  string `json:"uri"`
					Href string `json:"href"`
				} `json:"links"`
				Configuration struct {
				} `json:"configuration"`
				Quantity int `json:"quantity"`
			} `json:"_addtocartform"`
		} `json:"_element"`
	} `json:"_items"`
}
