#!/bin/sh

PROJECT=project-8979611309653332047
DATABASE_INSTANCE=lszt
ENTRY_POINT=CardPaymentsStripe
RUNTIME=go120

gcloud config set project $PROJECT
gcloud functions deploy cardPaymentsStripe \
  --entry-point $ENTRY_POINT \
  --trigger-event providers/google.firebase.database/eventTypes/ref.create \
  --trigger-resource projects/_/instances/$DATABASE_INSTANCE/refs/card-payments/{pushId} \
  --runtime $RUNTIME \
  --env-vars-file .env.prod


gcloud functions deploy StripeWebhook --runtime go120 --trigger-http --allow-unauthenticated --env-vars-file .env.prod
