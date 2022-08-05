# Traefik Forward Cookie Plugin
Because of the current(05-08-2022) implementation of the forward-auth middleware within traefik
you currently are unable to receive cookies set by the endpoint for forward-auth.

This plugin provides the necessary basis for forwarding an http-request and handling the cookie
that is provided from the endpoint.

## Configuration
In order to get this up and running you are required to configure the following fields:
```
testData:
  addr: "my.addr.com"
  cookies:
    - "response-cookie"
  headers:
    - "request-header"
  parameters:
    - "query-param"
```
* `addr`: the target address to send HTTP/HTTPs request
* `cookies`: list of all the cookies that you want to receive to the client
* `headers`: list of headers that may be required by the HTTP/HTTPs request
* `parameters`: list of all of the parameters that may be required by the HTTP/HTTPs request
