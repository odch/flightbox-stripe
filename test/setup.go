package test

import (
	"log"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhookendpoint"
)

func SetupWebhook(config *Config) error {

	stripe.Key = config.StripeSecret

	params := &stripe.WebhookEndpointParams{
		EnabledEvents: []*string{
			stripe.String("payment_intent.succeeded"),
			stripe.String("payment_intent.canceled"),
			stripe.String("payment_intent.payment_failed"),
		},
		URL: stripe.String(config.WebHookUrl),
	}
	result, err := webhookendpoint.New(params)
	log.Println(result.Secret)
	return err
}
