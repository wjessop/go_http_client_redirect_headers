# Go HTTP missing headers test

It seems the Go net/http client silently follows 30x redirects, but drops headers from the original request when doing so. This repo is for testing that out.

