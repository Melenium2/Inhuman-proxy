# Inhuman-reverse-proxy

This service needed for forwarding request to next services. Url of services need to set up inside environment variables.

This service additionally supports load balancing by fiber.

## Proxy

For each request this service add available proxy link to special header and next services can use this proxy for requesting
to external resources.