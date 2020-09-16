package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"pokemon_center_aco/models"
	"strings"

	"go.uber.org/zap"
)

const (
	RequestTypeGET  = "GET"
	RequestTypePOST = "POST"

	HomepageURL         = "https://www.pokemoncenter.com/"
	GetProductInfoURL   = "https://www.pokemoncenter.com/tpci-ecommweb-api/product?format=zoom.nodatalinks"
	AtcURL              = "https://www.pokemoncenter.com/tpci-ecommweb-api/cart?type=product&format=zoom.nodatalinks"
	SubmitProfileURL    = "https://www.pokemoncenter.com/tpci-ecommweb-api/address?format=zoom.nodatalinks"
	GetPaymentKeyuRL    = "https://www.pokemoncenter.com/tpci-ecommweb-api/payment/key?microform=true&locale=en-US"
	GetSubmitPaymentURL = "https://www.pokemoncenter.com/tpci-ecommweb-api/payment?microform=true&format=zoom.nodatalinks"
	SubmitOrderURL      = "https://www.pokemoncenter.com/tpci-ecommweb-api/order?format=zoom.nodatalinks"
)

type Client struct {
	client *http.Client
	logger *zap.Logger
	cookie string
	token  string
}

func InitializeClient(logger *zap.Logger) (Client, error) {
	c := Client{
		client: &http.Client{},
		logger: logger,
	}

	cookie, token, err := c.GetCookies()
	if err != nil {
		return c, err
	}

	c.cookie = cookie
	c.token = token

	return c, nil
}

func (c *Client) GetCookies() (string, string, error) {
	req, err := c.createRequest(RequestTypeGET, HomepageURL, "", false)
	if err != nil {
		c.logger.Error("GetCookies - httpNewRequest", zap.Error(err))
		return "", "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("GetCookies - sendRequest", zap.Error(err))
		return "", "", err
	}

	cookie := resp.Header.Get("set-cookie")
	authStr := getStringInBetweenInclusive(cookie, "{", "}")

	var auth models.CookieAuth
	err = json.Unmarshal([]byte(authStr), &auth)
	if err != nil {
		c.logger.Error("GetCookies - unmarshal", zap.Error(err))
		return "", "", err
	}

	return cookie, auth.AccessToken, nil
}

func (c *Client) GetProductAvailability(sku string) (bool, models.Product, error) {
	requestBody := fmt.Sprintf(`{"productSku":"%s"}`, sku)
	req, err := c.createRequest(RequestTypePOST, GetProductInfoURL, requestBody, true)
	if err != nil {
		c.logger.Error("GetProductAvailability - httpNewRequest", zap.Error(err))
		return false, models.Product{}, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("GetProductAvailability - sendRequest", zap.Error(err))
		return false, models.Product{}, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("GetProductAvailability - readResponseBody", zap.Error(err))
		return false, models.Product{}, err
	}

	var product models.Product
	err = json.Unmarshal(bodyText, &product)
	if err != nil {
		c.logger.Error("GetProductAvailability - unmarshal", zap.Error(err))
		return false, models.Product{}, err
	}

	a := product.Availability[0].State
	if a == PRODUCT_AVAILABLE {
		return true, product, nil
	}

	return false, product, nil
}

func (c *Client) AddToCart(uri string, quantity int) (bool, error) {
	requestBody := fmt.Sprintf(`{"productURI":"%s","quantity":%d}`, uri, quantity)
	req, err := c.createRequest(RequestTypePOST, GetProductInfoURL, requestBody, true)
	if err != nil {
		c.logger.Error("AddToCart - httpNewRequest", zap.Error(err))
		return false, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("AddToCart - sendRequest", zap.Error(err))
		return false, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("AddToCart - readResponseBody", zap.Error(err))
		return false, err
	}

	var atc models.ATC
	err = json.Unmarshal(bodyText, &atc)
	if err != nil {
		c.logger.Error("AddToCart - unmarshal", zap.Error(err))
		return false, err
	}

	messages := atc.Messages
	if len(messages) == 0 {
		return true, nil
	}

	return false, nil
}

func (c *Client) SubmitBillingShippingInfo(profile models.Profile) error {
	requestBody, err := json.Marshal(profile)
	if err != nil {
		c.logger.Error("SubmitBillingShippingInfo - marshalProfile", zap.Error(err))
		return err
	}

	req, err := c.createRequest(RequestTypePOST, SubmitProfileURL, string(requestBody), true)
	if err != nil {
		c.logger.Error("SubmitBillingShippingInfo - httpNewRequest", zap.Error(err))
		return err
	}

	_, err = c.client.Do(req)
	if err != nil {
		c.logger.Error("SubmitBillingShippingInfo - sendRequest", zap.Error(err))
		return err
	}

	return nil
}

func (c *Client) GetPaymentKey() (string, error) {
	req, err := c.createRequest(RequestTypeGET, GetPaymentKeyuRL, "", true)
	if err != nil {
		c.logger.Error("GetPaymentKey - httpNewRequest", zap.Error(err))
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("GetPaymentKey - sendRequest", zap.Error(err))
		return "", err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("GetPaymentKey - readResponseBody", zap.Error(err))
		return "", err
	}

	var paymentKey models.PaymentKey
	err = json.Unmarshal(bodyText, &paymentKey)
	if err != nil {
		c.logger.Error("GetPaymentKey - unmarshal", zap.Error(err))
		return "", err
	}

	return paymentKey.KeyID, nil
}

func (c *Client) GetSubmitOrderURL(body models.PaymentURLRequest) (string, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		c.logger.Error("GetSubmitPaymentURL - marshalProfile", zap.Error(err))
		return "", err
	}

	req, err := c.createRequest(RequestTypePOST, GetSubmitPaymentURL, string(requestBody), true)
	if err != nil {
		c.logger.Error("GetSubmitPaymentURL - httpNewRequest", zap.Error(err))
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("GetSubmitPaymentURL - sendRequest", zap.Error(err))
		return "", err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("GetSubmitPaymentURL - readResponseBody", zap.Error(err))
		return "", err
	}

	var url models.PaymentURLResponse
	err = json.Unmarshal(bodyText, &url)
	if err != nil {
		c.logger.Error("GetSubmitPaymentURL - unmarshal", zap.Error(err))
		return "", err
	}

	orderLink, err := convertPaymentLinkToOrderLink(url.Self.URI)
	if err != nil {
		c.logger.Error("GetSubmitPaymentURL - convertPaymentLink", zap.Error(err), zap.String("payment link", url.Self.URI))
		return "", err
	}

	return orderLink, nil
}

func (c *Client) SubmitOrder(url string) (string, error) {
	requestBody := fmt.Sprintf(`{"purchaseForm":"%s"}`, url)
	req, err := c.createRequest(RequestTypePOST, SubmitOrderURL, requestBody, true)
	if err != nil {
		c.logger.Error("SubmitOrder - httpNewRequest", zap.Error(err))
		return "", err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.Error("SubmitOrder - sendRequest", zap.Error(err))
		return "", err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("SubmitOrder - readResponseBody", zap.Error(err))
		return "", err
	}

	var orderResponse models.SubmitOrderResponse
	err = json.Unmarshal(bodyText, &orderResponse)
	if err != nil {
		c.logger.Error("SubmitOrder - unmarshal", zap.Error(err))
		return "", err
	}

	messages := orderResponse.Messages
	if len(messages) != 0 {
		err = errors.New(messages[0].DebugMessage)
		c.logger.Error("SubmitOrder - failed to submit order", zap.Error(err))
		return "", err
	}

	return orderResponse.PurchaseNumber, nil
}

func (c *Client) createRequest(requestType, url, body string, withCookies bool) (*http.Request, error) {
	var req *http.Request
	var err error

	if requestType == RequestTypeGET {
		req, err = http.NewRequest(requestType, url, nil)
	} else {
		req, err = http.NewRequest(requestType, url, strings.NewReader(body))
	}
	if err != nil {
		return nil, err
	}

	return c.setHeaders(req, withCookies), nil
}

func (c *Client) setHeaders(req *http.Request, withCookies bool) *http.Request {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:80.0) Gecko/20100101 Firefox/80.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-Store-Scope", "pokemon")
	req.Header.Set("Origin", "https://www.pokemoncenter.com")

	if withCookies {
		req.Header.Set("Cookie", c.cookie)
		req.Header.Set("Authorization", fmt.Sprintf("bearer %s", c.token))
	}

	return req
}

// Order link formatting: "/paymentmethods/orders/pokemon/xxxxxx=" -> "/purchases/orders/pokemon/xxxxxx=/form"
func convertPaymentLinkToOrderLink(url string) (string, error) {
	prefixToRemove := "/paymentmethods/orders/pokemon/"
	prefixToAdd := "/purchases/orders/pokemon/"
	suffixToAdd := "/form"

	if !strings.Contains(url, prefixToRemove) {
		return url, errors.New("Submit payment URL has incorrect formatting")
	}

	return strings.Join(
		[]string{
			strings.Replace(url, prefixToRemove, prefixToAdd, -1),
			suffixToAdd,
		}, "",
	), nil
}
