package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"html"
	"log"
	"context"
	"crypto/tls"
	"flag"
	"io"
	"time"

	"golang.org/x/crypto/acme/autocert"
)

var dict map[string]Unit

type Unit struct {
	Version string `json:"version"`
	BeanID 	string `json:"beanID"`
}

// Utility

func checkError(message string, err error) {
	if err != nil {
		fmt.Println("ERROR: " + message)
		fmt.Println(err)
		fmt.Println("")
	}
}

func formatRequest(r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	checkError("couldnt format http request", err)
	log.Println(string(requestDump))
}

// HTTP Handlers

func unitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. unit handler here")
	// TODO get id, version, and beanID from request, then publish update msg on mqtt channel with that data
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. list handler here")
	// TODO convert dict to json, send that as response
}

func makeHTTPServer() *http.Server {
	router := &http.ServeMux{}
	//mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/unit/{beanID}", unitHandler)
	router.HandleFunc("/units/{labID}", listHandler)
	router.Handle("/lab/{labID}", http.StripPrefix("/lab/{labID}/", http.FileServer(http.Dir("./static"))))
	//router.PathPrefix("/lab/{labID}").Handler(http.FileServer(http.Dir("./static")))

	return &http.Server{
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
        IdleTimeout:  120 * time.Second,
        Handler:      router,
    }
}

func main() {
	dict = make(map[string]Unit)

	var m *autocert.Manager
	var server *http.Server

	hostPolicy = func(ctx context.Context, host string) error {
		allowedHost = "www.sbc.crabdance.com"
		if host == allowedHost {
			return nil
		}

		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	}

	certPath := "/etc/letsencrypt/live/sbc.crabdance.com/"
	m = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(certPath),
	}

	server = makeHTTPServer()
	server.Addr = ":443"
	server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

	go func() {
		log.Println("Starting server on %s\n", server.Addr)
		log.Fatal(server.ListenAndServeTLS("", ""))
	}()
}