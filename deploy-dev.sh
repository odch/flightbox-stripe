#!/bin/sh
set -e
PROJECT=mfgt-flights
DATABASE_INSTANCE=mfgt-flights
ENTRY_POINT=CardPaymentsStripe
RUNTIME=go120

gcloud config set project $PROJECT
gcloud functions deploy cardPaymentsStripe \
  --entry-point $ENTRY_POINT \
  --trigger-event providers/google.firebase.database/eventTypes/ref.create \
  --trigger-resource projects/_/instances/$DATABASE_INSTANCE/refs/card-payments/{pushId} \
  --runtime $RUNTIME \
  --env-vars-file .env.dev


gcloud functions deploy StripeWebhook --runtime go120 --trigger-http --allow-unauthenticated --env-vars-file .env.dev
