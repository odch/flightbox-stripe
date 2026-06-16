package p

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
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
	functions.HTTP("StripeWebhook", StripeWebhook)
	functions.CloudEvent("CardPaymentsStripe", cardPaymentsStripe)
}

func StripeWebhook(w http.ResponseWriter, req *http.Request) {
	// Without a configured signing secret, event signatures would be
	// verified against the public default value, so reject the request
	// instead of falling back to it.
	if config.WebHookSecret == "" || config.WebHookSecret == "not_configured" {
		log.Println("WEBHOOK_SECRET is not configured, rejecting request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

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

func cardPaymentsStripe(ctx context.Context, e event.Event) error {
	var rtdbEvent RTDBEvent
	if err := e.DataAs(&rtdbEvent); err != nil {
		return fmt.Errorf("event.DataAs: %w", err)
	}

	subject := e.Subject()
	log.Printf("Function triggered by change to: %v", subject)
	idx := strings.Split(subject, "/")
	id := idx[len(idx)-1]
	log.Printf("%+v", rtdbEvent)

	var err error
	if rtdbEvent.Delta.Method == "card" {
		err = test.TerminalPayment(config, id, rtdbEvent.Delta.Amount, &rtdbEvent.Delta.Email, rtdbEvent.Delta.Registration)
	} else {
		err = test.CheckoutPayment(config, id, rtdbEvent.Delta.Amount, rtdbEvent.Delta.Email, rtdbEvent.Delta.Registration, rtdbEvent.Delta.ArrivalReference, rtdbEvent.Delta.RefNr)
	}
	if err != nil {
		log.Println(err)
	}
	return nil
}
