#!/bin/sh
set -e

# Set the project id and database name here
PROJECT=XXX
DATABASE_INSTANCE=XXX

RUNTIME=go121

gcloud config set project $PROJECT

gcloud functions deploy cardPaymentsStripe \
  --entry-point CardPaymentsStripe \
  --trigger-event providers/google.firebase.database/eventTypes/ref.create \
  --trigger-resource projects/_/instances/$DATABASE_INSTANCE/refs/card-payments/{pushId} \
  --runtime $RUNTIME \
  --no-gen2 \
  --env-vars-file .env.dev

gcloud functions deploy StripeWebhook \
  --allow-unauthenticated \
  --trigger-http \
  --runtime $RUNTIME \
  --no-gen2 \
  --env-vars-file .env.dev
gcloud functions add-iam-policy-binding StripeWebhook \
  --region=us-central1 \
  --member=allUsers \
  --role=roles/cloudfunctions.invoker