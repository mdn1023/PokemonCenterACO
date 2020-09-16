package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"pokemon_center_aco/models"
	"strings"

	"go.uber.org/zap"
)

const (
	PRODUCT_AVAILABLE   = "AVAILABLE"
	PRODUCT_UNAVAILABLE = "NOT_AVAILABLE"
)

func GetCookies(logger *zap.Logger) (string, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.pokemoncenter.com/", nil)
	if err != nil {
		logger.Error("GetCookies - httpNewRequest", zap.Error(err))
		return "", "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("X-Store-Scope", "pokemon")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("GetCookies - sendRequest", zap.Error(err))
		return "", "", err
	}

	cookie := resp.Header.Get("set-cookie")
	authStr := getStringInBetweenInclusive(cookie, "{", "}")

	var auth models.CookieAuth
	err = json.Unmarshal([]byte(authStr), &auth)
	if err != nil {
		logger.Error("GetCookies - unmarshal", zap.Error(err))
		return "", "", err
	}

	return cookie, auth.AccessToken, nil
}

func GetProductAvailability(sku, cookie, bearerToken string, logger *zap.Logger) (bool, models.Product, error) {
	client := &http.Client{}

	requestBody := fmt.Sprintf(`{"productSku":"%s"}`, sku)
	var data = strings.NewReader(requestBody)
	req, err := http.NewRequest("POST", "https://www.pokemoncenter.com/tpci-ecommweb-api/product?format=zoom.nodatalinks", data)
	if err != nil {
		logger.Error("GetProductAvailability - httpNewRequest", zap.Error(err))
		return false, models.Product{}, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Store-Scope", "pokemon")
	req.Header.Set("Origin", "https://www.pokemoncenter.com")
	req.Header.Set("Connection", "keep-alive")

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearerToken))

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("GetProductAvailability - sendRequest", zap.Error(err))
		return false, models.Product{}, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("GetProductAvailability - readResponseBody", zap.Error(err))
		return false, models.Product{}, err
	}

	var product models.Product
	err = json.Unmarshal(bodyText, &product)
	if err != nil {
		logger.Error("GetProductAvailability - unmarshal", zap.Error(err))
		return false, models.Product{}, err
	}

	a := product.Availability[0].State
	if a == PRODUCT_AVAILABLE {
		return true, product, nil
	}

	return false, product, nil
}

func AddToCart(uri, cookie, bearerToken string, quantity int, logger *zap.Logger) (bool, error) {
	client := &http.Client{}

	requestBody := fmt.Sprintf(`{"productURI":"%s","quantity":%d}`, uri, quantity)
	var data = strings.NewReader(requestBody)
	req, err := http.NewRequest("POST", "https://www.pokemoncenter.com/tpci-ecommweb-api/cart?type=product&format=zoom.nodatalinks", data)
	if err != nil {
		logger.Error("AddToCart - httpNewRequest", zap.Error(err))
		return false, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Store-Scope", "pokemon")
	req.Header.Set("Origin", "https://www.pokemoncenter.com")
	req.Header.Set("Connection", "keep-alive")

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearerToken))

	resp, err := client.Do(req)
	if err != nil {
		logger.Error("AddToCart - sendRequest", zap.Error(err))
		return false, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error("AddToCart - readResponseBody", zap.Error(err))
		return false, err
	}

	var atc models.ATC
	err = json.Unmarshal(bodyText, &atc)
	if err != nil {
		logger.Error("AddToCart - unmarshal", zap.Error(err))
		return false, err
	}

	messages := atc.Messages
	if len(messages) == 0 {
		return true, nil
	}

	return false, nil
}

func SubmitBillingShippingInfo(profile models.Profile, cookie, bearerToken string, logger zap.Logger) {
	client := &http.Client{}

	marshalledBody, err := json.Marshal(profile)
	if err != nil {
		fmt.Println("error:", err)
	}

	var data = strings.NewReader(string(marshalledBody))
	req, err := http.NewRequest("POST", "https://www.pokemoncenter.com/tpci-ecommweb-api/address?format=zoom.nodatalinks", data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Store-Scope", "pokemon")
	req.Header.Set("Origin", "https://www.pokemoncenter.com")
	req.Header.Set("Connection", "keep-alive")

	req.Header.Set("Cookie", cookie)
	req.Header.Set("Authorization", fmt.Sprintf("bearer %s", bearerToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", bodyText)
}

func getStringInBetweenInclusive(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start) - 1
	e := strings.Index(str, end) + 1
	if e == -1 {
		return
	}
	return str[s:e]
}
