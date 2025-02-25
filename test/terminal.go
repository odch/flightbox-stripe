package test

import (
	"log"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/testhelpers/terminal/reader"
)

func Simulate(config *Config) error {
	//var visa string = "4242424242424242"

	stripe.Key = config.StripeSecret
	params := &stripe.TestHelpersTerminalReaderPresentPaymentMethodParams{
		// CardPresent: &stripe.TestHelpersTerminalReaderPresentPaymentMethodCardPresentParams{
		// 	Number: &visa,
		// },
	}
	result, err := reader.PresentPaymentMethod("tmr_FFCq9gvxhcRejj", params)
	log.Println(result)
	return err
}
