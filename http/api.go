package http

import (
	"pokemon_center_aco/models"

	"go.uber.org/zap"
)

type API interface {
	GetCookies(logger *zap.Logger) (string, string, error)
	GetProductAvailability(sku, cookie, bearerToken string, logger *zap.Logger) (bool, models.Product, error)
	AddToCart(uri, cookie, bearerToken string, quantity int, logger *zap.Logger) (bool, error)
	SubmitBillingShippingInfo(profile models.Profile, cookie, bearerToken string, logger zap.Logger)
}
