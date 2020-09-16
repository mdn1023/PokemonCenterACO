package models

// CookieAuth respresents the json struct of the cookie auth generated in the response header when making a GET request to the PokemonCenter homepage
type CookieAuth struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int      `json:"expires_in"`
	Scope       string   `json:"scope"`
	Role        string   `json:"role"`
	Roles       []string `json:"roles"`
}
