package main

import "github.com/odch/flightbox/functions-go/stripe-terminal/test"

func main() {
	config, err := test.LoadConfig()
	if err != nil {
		panic(err)
	}
	if err := test.SetupWebhook(config); err != nil {
		panic(err)
	}
}
