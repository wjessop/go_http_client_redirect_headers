# Go HTTP missing headers test

It seems the Go net/http client silently follows 30x redirects, but drops headers from the original request when doing so. This repo is for testing that out.

The program makes starts up a Go server that has two handlers. The first redirects to the second. Both dump out the headers they have received. the program then makes two requests to the server. The first is a "vanilla" request that doesn't specify a redirect handler. The second makes a request, but this time specifies a redirect handler that copies the header values between requests.

There's been some [discussion on this previously](http://grokbase.com/t/gg/golang-nuts/136syegzsc/go-nuts-net-http-redirection-and-headers) but no resolution, wether this is correct or incorrect behaviour that I can find.

## The output

You can generate this output yourself by running:

    go run main.go

However here it is to save you the trouble, I've split it into the two requests to make it easier to parse

### Request one

This is the "vanilla" request with no redirect policy function. You can see that the client makes the request to the /start handler and the next it hears back is the response from the /redirectdestination handler. The client lost the two header values that were specified.

    * Client making request
    Server got request: GET /start
    Headers:
    	User-Agent: Go 1.1 package http
    	Content-Type: application/json
    	X-Foo: bar
    	Accept-Encoding: gzip
    Redirecting to /redirectdestination

    Server got request: GET /redirectdestination
    Headers:
    	User-Agent: Go 1.1 package http
    	Referer: http://localhost:8080/start
    	Accept-Encoding: gzip
    * Client response: 200 OK

### Request two

Now we make another request, this time we specify a redirect policy function to copy the headers into subsequent requests. You can see that the client gets control back after the redirect (the lines starting with *)

    * Client making request
    Server got request: GET /start
    Headers:
    	User-Agent: Go 1.1 package http
    	Content-Type: application/json
    	X-Foo: bar
    	Accept-Encoding: gzip
    Redirecting to /redirectdestination

    * Client handling redirect
    * Previous requests: 1
    * 	1: GET http://localhost:8080/start
    * Copying headers
    * Following redirect
    Server got request: GET /redirectdestination
    Headers:
    	Content-Type: application/json
    	Referer: http://localhost:8080/start
    	X-Foo: bar
    	Accept-Encoding: gzip
    	User-Agent: Go 1.1 package http
    * Client response: 200 OK

## A warning

Don't use this redirect policy func in production code, there are issues with it, it's just there to demonstrate this header issue
