Websocket API Gateway

This application uses Gin and Gorilla websocket to implement a gateway for subscribing
to the signals of the test account of the Bitmex exchange, when connected to which
it is possible to receive notifications about changes in the table of financial
instruments (quotes) of the test environment of the Bitmex exchange.

You can run service by command
"go run main.go"

Test as a client using websocat by command
"websocat ws://127.0.0.1:8844/bitmex"


