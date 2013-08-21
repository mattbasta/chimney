package main

import (
	"flag"
	"fmt"
    "io"
    "io/ioutil"
	"log"
	"net/http"
	"time"
)

var api = flag.String("api", "", "the api hostname to proxy")
var port = flag.Int("port", 8080, "the port to listen on")

type Proxy struct {
	client http.Client
}

func (self *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u := r.URL
	u.Scheme = "https"
	u.Host = *api

	nreq, err := http.NewRequest(r.Method, u.String(), r.Body)
	if err != nil {
		panic(err)
	}

	// Write request headers to the server
	for k, v := range r.Header {
		nreq.Header.Set(k, v[0])
	}

    resp, err := self.client.Do(nreq)
    if err != nil {
        panic(err)
    }
    defer io.Copy(ioutil.Discard, resp.Body)
    defer resp.Body.Close()
	// go func() {
	// 	fmt.Printf(r.URL.String() + "\n")
	// }()

	// Write response headers to the client
	for k, v := range resp.Header {
		w.Header().Set(k, v[0])
	}
    w.WriteHeader(resp.StatusCode)

    // Read the response out from the server.
    contents, err := ioutil.ReadAll(resp.Body)
    if err != nil && err != io.EOF {
        panic(err)
    }
    // Write the response out to the client.
    if _, err := w.Write(contents); err != nil {
        panic(err)
    }

}

func NewProxy() Proxy {
	return Proxy{http.Client{}}
}

func main() {
	flag.Parse()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      new(Proxy),
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}

	log.Fatal(server.ListenAndServe())

}
