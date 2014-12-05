package main

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	serv_address   string = "localhost:8080"
	http_start_url string = "http://localhost:8080/start"
)

func main() {
	// A HTTP server to redirect and log headers
	go http_server()

	fmt.Println(`
===================================================

First we make a request where we add a
header. The header won't make it past the
redirect

===================================================`)

	client := &http.Client{}
	makeRequest(client)

	fmt.Println(`
===================================================

Now we make a request, but this time we
specify a redirectPolicyFunc to copy over
header values

===================================================`)

	client = &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	makeRequest(client)
}

func makeRequest(client *http.Client) {
	req, err := http.NewRequest("GET", http_start_url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-foo", "bar")
	fmt.Println("* Client making request")
	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	// fmt.Println("----------------------")
	fmt.Println("* Client response:", resp.Status)
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	// fmt.Println("**********************")
	fmt.Println("* Client handling redirect")
	fmt.Printf("* Previous requests: %d\n", len(via))
	for i, req := range via {
		fmt.Printf("* \t%d: %s %s", i+1, req.Method, req.URL)
		fmt.Println()
	}
	fmt.Println("* Copying headers redirect")
	// GET &{GET http://localhost:8080/start HTTP/1.1 %!s(int=1) %!s(int=1) map[Content-Type:[application/json]] <nil> %!s(int64=0) [] %!s(bool=false) localhost:8080 map[] map[] %!s(*multipart.Form=<nil>) map[]   %!s(*tls.ConnectionState=<nil>)}
	// Copy the headers out of the last request we made
	for key, vals := range via[len(via)-1].Header {
		for _, val := range vals {
			req.Header.Add(key, val)
		}

	}

	fmt.Println("* Following redirect")
	// fmt.Println("**********************")
	return nil
}

func http_server() {
	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("----------------------")
		fmt.Printf("Server got request: %s %s\n", r.Method, r.RequestURI)
		fmt.Println("Headers:")

		for key, values := range r.Header {
			fmt.Printf("\t%s: %s\n", key, strings.Join(values, ", "))
		}

		fmt.Println("Redirecting to /redirectdestination\n")
		http.Redirect(w, r, "/redirectdestination", 301)
	})

	http.HandleFunc("/redirectdestination", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("----------------------")
		fmt.Printf("Server got request: %s %s\n", r.Method, r.RequestURI)
		fmt.Println("Headers:")

		for key, values := range r.Header {
			fmt.Printf("\t%s: %s\n", key, strings.Join(values, ", "))
		}
	})

	panic(http.ListenAndServe(serv_address, nil))
}
