package test

import (
	"context"
	"log"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/paymentlink"
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

func PaylinkPayment(config *Config, id string, amount int64, receiptEMail *string, registration string) error {
	stripe.Key = config.StripeSecret

	qty := int64(1)
	pid := "price_1QwL8RQ837JfjWrVQLa3808e" // roli test price
	//pid := "prod_Pn8NbuL38Z8D4n" // philipp test product
	linkParams := stripe.PaymentLinkParams{
		LineItems: []*stripe.PaymentLinkLineItemParams{
			{
				Price:    &pid,
				Quantity: &qty,
			},
		},
	}
	linkParams.AddMetadata("external_id", id)
	linkParams.AddMetadata("registration", registration)

	result, err := paymentlink.New(&linkParams)

	if err != nil {
		return err
	}
	log.Println(result)
	ctx := context.TODO()
	UpdateData(config, ctx, id, result.URL)
	UpdateStatus(config, ctx, id, "pending")

	return nil
}
