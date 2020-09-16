package main

import (
	"fmt"
	"pokemon_center_aco/http"
	"pokemon_center_aco/models"
	"time"

	"go.uber.org/zap"
)

const (
	SHINING_LEGENDS_ETB_SKU = "290-80319" // Out of stock
	REBEL_CLASH_ETB_SKU     = "173-80700" // In stock
	REBEL_CLASH_BOOSTER_SKU = "173-80682" // Out of stock

	SHINING_LEGENDS_ETB_ATC_URI = "/carts/items/pokemon/qgqvhkjsheyc2obqgmyts=/form"
	REBEL_CLASH_ETB_ATC_URI     = "/carts/items/pokemon/qgqvhkjrg4zs2obqg4yda=/form"
	REBEL_CLASH_BOOSTER_ATC_URI = "/carts/items/pokemon/qgqvhkjrg4zs2obqgy4de=/form"

	MONITOR_DELAY = 5000
	QUANTITY      = 1
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any

	// Start session, grab cookies
	logger.Info("Initializing session...grabbing cookie...")
	client, err := http.InitializeClient(logger)
	if err != nil {
		logger.Fatal("Could not retrieve session cookie", zap.Error(err))
	}
	logger.Info("Session created!")

	// Check availability of the product by SKU on a delay
	product := models.Product{}
	for {
		logger.Info("Checking product availability...")
		available, p, err := client.GetProductAvailability(REBEL_CLASH_ETB_SKU)
		if err != nil {
			logger.Fatal("Check product availability failed", zap.Error(err))
		}

		if available {
			logger.Info("Product available!")
			product = p
			break
		}

		logger.Info("Product unavailable! Retrying....")
		time.Sleep(MONITOR_DELAY * time.Millisecond)
	}

	// Once product has been found and is available, parse out the ATC url and attempt to add item to cart
	atcURL := product.Items[0].Element[0].AddToCartForm[0].Self.URI
	for {
		logger.Info("Adding product to cart...", zap.String("ATC URL", atcURL))
		carted, err := client.AddToCart(atcURL, QUANTITY)

		if err != nil {
			logger.Fatal("Add to cart failed, error", zap.Error(err))
		}

		if carted {
			logger.Info("Product added to cart!")
			break
		}

		logger.Info("Couldn't add to cart, OOS! Retrying....")
		time.Sleep(MONITOR_DELAY * time.Millisecond)
	}

	// Construct the billing/shipping info and send it to the server
	profile := models.Profile{
		Billing: models.Address{
			FirstName:    "Michael",
			LastName:     "Nguyen",
			Country:      "US",
			City:         "Rockville",
			PhoneNumber:  "(240) 447-6251",
			PostalCode:   "20850",
			State:        "MD",
			AddressLine1: "3 McCormick Court",
			AddressLine2: "Apt 4",
		},
		Shipping: models.Address{
			FirstName:    "Michael",
			LastName:     "Nguyen",
			Country:      "US",
			City:         "Rockville",
			PhoneNumber:  "(240) 447-6251",
			PostalCode:   "20850",
			State:        "MD",
			AddressLine1: "3 McCormick Court",
			AddressLine2: "Apt 4",
		},
	}

	logger.Info("Submitting shipping and billing info...")
	err = client.SubmitBillingShippingInfo(profile)
	if err != nil {
		logger.Fatal("Error submitting shipping and billing info", zap.Error(err))
	}
	logger.Info("Shipping and billing info successfully submitted!")

	// Retrieve the session's payment key ID, used in conjunction with payment data and payment token to generate a final checkout url
	logger.Info("Retrieving payment key ID...")
	paymentKeyID, err := client.GetPaymentKey()
	if err != nil {
		logger.Fatal("Error retrieving session payment key ID", zap.Error(err))
	}
	logger.Info("Successfully retrieved payment key ID!")
	fmt.Println(paymentKeyID)

	// TODO: Grab the payment token hidden through the source code javascript obfuscation
	// TODO: All code below will not work until payment token can be retrieved

	// Send the payment token, payment key ID, and payment info through 'GetSubmitPaymentURL' to get the final submit order URL
	paymentInfo := models.PaymentInfo{
		Name:             "John Doe",
		CreditCardType:   "VISA",
		CreditCardNumber: "4111111111111111",
		Cvv:              "015",
		ExpiryMonth:      "09",
		ExpiryYear:       "2025",
	}

	body := models.PaymentURLRequest{
		PaymentDisplay: fmt.Sprintf("%s %s/%s", paymentInfo.CreditCardType, paymentInfo.ExpiryMonth, paymentInfo.ExpiryYear),
		PaymentKey:     paymentKeyID,
		PaymentToken:   "FakePaymentToken - Received from above request",
	}

	// Retrieve the submit order url
	logger.Info("Retrieving submit order URL...")
	url, err := client.GetSubmitOrderURL(body)
	if err != nil {
		logger.Fatal("Error retrieving submit order URL", zap.Error(err))
	}
	logger.Info("Successfully retrieved submit order URL!", zap.String("SubmitOrder URL: ", url))

	// Submit the order!
	logger.Info("Submitting order...")
	orderID, err := client.SubmitOrder(url)
	if err != nil {
		logger.Fatal("Error submitting order :(", zap.Error(err))
	}
	logger.Info("YAY! Order successfully submitted!", zap.String("Order ID: ", orderID))
}

//{"purchaseForm":"/purchases/orders/pokemon/gbqtomdcmm3gcljymu4diljumy4wcllcgy4taljtgqydizjsgnstgndfgm=/form"}

// 1C2LATBTEUS510BWUJHIGS89P9HQEJ4V3SVIJSXKTYIBMGMFXMZA5F5FFE14854E

// 1C1IMOQYYN4G9J6OBKYFKYT98ANZXZLPUO4NGJ4BH5ACPO3AK4BG5F5FFE471903

/**

t.Base64 = {
            encode: function(e) {
                var t = u(e) ? l(e) : e;
                if (!(t instanceof ArrayBuffer)) throw new TypeError("Input should be of type String or ArrayBuffer");
                var n = new Uint8Array(t);
                return window.btoa(String.fromCharCode.apply(null, n))
            }
		}

*/

/*
	alg: "RSA-OAEP",
	enc: "A256GCM"

*/
