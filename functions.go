package p

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"cloud.google.com/go/functions/metadata"
	"github.com/odch/flightbox/functions-go/stripe-terminal/test"
	"github.com/stripe/stripe-go/v74/webhook"
)

var config *test.Config

func init() {
	var err error
	config, err = test.LoadConfig()
	if err != nil {
		panic(err)
	}
	//functions.HTTP("StripeWebhook", StripeWebhook)
}

func StripeWebhook(w http.ResponseWriter, req *http.Request) {
	// Protects against a malicious client streaming us an endless request
	// body
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Pass the request body & Stripe-Signature header to ConstructEvent, along with the webhook signing key
	event, err := webhook.ConstructEvent(body, req.Header.Get("Stripe-Signature"), config.WebHookSecret)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		fmt.Fprintf(w, "%v", err)
		return
	}
	ctx := context.Background()

	fmt.Fprintf(w, "Received signed event: %v", event)

	id := event.GetObjectValue("metadata", "external_id")
	if id != "" {
		if event.Type == "payment_intent.succeeded" {
			test.UpdateStatus(config, ctx, id, "success")

		} else if event.Type == "payment_intent.canceled" ||
			event.Type == "payment_intent.payment_failed" {
			test.UpdateStatus(config, ctx, id, "failure")
		}
	}
}

// RTDBEvent is the payload of a RTDB event.
type RTDBEvent struct {
	Data  interface{} `json:"data"`
	Delta struct {
		Amount           int64  `json:"amount"` // cents
		ArrivalReference string `json:"arrivalReference"`
		RefNr            string `json:"refNr"`
		Currency         string `json:"currency"`
		Email            string `json:"email"`
		Registration     string `json:"immatriculation"`
		Method           string `json:"method"`
	} `json:"delta"`
}

// {Data:<nil> Delta:map[amount:16 arrivalReference:-N_JibawWc-1GbZbxMzh currency:CHF status:pending timestamp:1.689345456348e+12]}
func CardPaymentsStripe(ctx context.Context, e RTDBEvent) error {
	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %w", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)
	log.Printf("%s - %s - %s - %s", meta.Resource.Name, meta.Resource.Service, meta.Resource.Type, meta.Resource.RawPath)
	idx := strings.Split(meta.Resource.RawPath, "/")
	id := idx[len(idx)-1]
	log.Printf("%+v", e)
	if e.Delta.Method == "card" {
		err = test.TerminalPayment(config, id, e.Delta.Amount, &e.Delta.Email, e.Delta.Registration)
	} else {
		err = test.CheckoutPayment(config, id, e.Delta.Amount, e.Delta.Email, e.Delta.Registration, e.Delta.ArrivalReference, e.Delta.RefNr)
	}
	if err != nil {
		log.Println(err)
	}
	return nil
}
