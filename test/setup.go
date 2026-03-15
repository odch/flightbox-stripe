package test

import (
	"log"
	"strings"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhookendpoint"
)

func SetupWebhook(config *Config) error {
	log.Printf("Setting up webhook for URL: %s", config.WebHookUrl)

	stripe.Key = config.StripeSecret

	// List existing webhooks
	params := &stripe.WebhookEndpointListParams{}
	iter := webhookendpoint.List(params)

	for iter.Next() {
		we := iter.WebhookEndpoint()
		if we.URL == config.WebHookUrl {
			log.Printf("Webhook already exists: %s", we.ID)
			return nil
		}
		if strings.Contains(we.URL, "StripeWebhook") {
			log.Printf("Updating webhook %s URL from %s to %s", we.ID, we.URL, config.WebHookUrl)
			updateParams := &stripe.WebhookEndpointParams{
				URL: stripe.String(config.WebHookUrl),
			}
			_, err := webhookendpoint.Update(we.ID, updateParams)
			if err != nil {
				return err
			}
			log.Printf("Webhook updated successfully")
			return nil
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	// Create new webhook with correct API version
	createParams := &stripe.WebhookEndpointParams{
		EnabledEvents: []*string{
			stripe.String("payment_intent.succeeded"),
			stripe.String("payment_intent.canceled"),
			stripe.String("payment_intent.payment_failed"),
		},
		URL:        stripe.String(config.WebHookUrl),
		APIVersion: stripe.String("2022-11-15"),
	}
	result, err := webhookendpoint.New(createParams)
	if err != nil {
		return err
	}

	log.Printf("Created new webhook: %s", result.ID)
	log.Println(result.Secret)
	return nil
}
