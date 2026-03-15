# Stripe Restricted Key Permissions

The following permissions are required for the restricted API key for online checkouts (external terminal requires more
permissions):

| Permission             | Level |
|------------------------|-------|
| Core - Payment Intents | Write |
| Checkout               | Write |
| Webhook                | Write |

All other permissions can remain at **None**.
