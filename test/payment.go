package test

import (
	"context"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/terminal/reader"
)

func TerminalPayment(config *Config, id string, amount int64, receiptEMail *string, registration string) error {

	stripe.Key = config.StripeSecret

	intentParams := &stripe.PaymentIntentParams{
		Currency:           stripe.String(string(stripe.CurrencyCHF)),
		PaymentMethodTypes: []*string{stripe.String("card_present")},
		CaptureMethod:      stripe.String(string(stripe.PaymentIntentCaptureMethodAutomatic)),
		Amount:             stripe.Int64(amount),
		ReceiptEmail:       receiptEMail,
	}
	intentParams.AddMetadata("external_id", id)
	intentParams.AddMetadata("registration", registration)

	intentResult, err := paymentintent.New(intentParams)
	if err != nil {
		return err
	}
	log.Println(intentResult)

	paymentParams := &stripe.TerminalReaderProcessPaymentIntentParams{
		PaymentIntent: stripe.String(intentResult.ID),
	}
	paymentResult, err := reader.ProcessPaymentIntent(config.TerminalId, paymentParams)

	if err != nil {
		return err
	}
	log.Println(paymentResult)

	return nil
}

func CheckoutPayment(config *Config, id string, amount int64, receiptEMail string, registration string, arrivalReference string, refNr string) error {
	stripe.Key = config.StripeSecret

	sessionParams := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency:   stripe.String(string(stripe.CurrencyCHF)),
					UnitAmount: stripe.Int64(amount),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Landetaxe"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String("payment"),
		SuccessURL: stripe.String(fmt.Sprintf("%s/#/arrival/%s/payment?success=true", config.ReturnBaseUrl, arrivalReference)),
		CancelURL:  stripe.String(fmt.Sprintf("%s/#/arrival/%s/payment?cancel=true", config.ReturnBaseUrl, arrivalReference)),
		Params: stripe.Params{
			Metadata: map[string]string{
				"registration": registration,
				"arrival":      arrivalReference,
				"refNr":        refNr,
			},
		},
		CustomerEmail: stripe.String(receiptEMail),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			ReceiptEmail: stripe.String(receiptEMail),
			Metadata: map[string]string{
				"registration": registration,
				"arrival":      arrivalReference,
				"refNr":        refNr,
			},
		},
	}

	s, err := session.New(sessionParams)
	if err != nil {
		log.Fatalf("Failed to create Checkout session: %v\n", err)
	}

	fmt.Println("Checkout Session URL:", s.URL)

	ctx := context.TODO()

	UpdateData(config, ctx, id, s.URL)
	UpdateStatus(config, ctx, id, "pending")

	return nil
}
