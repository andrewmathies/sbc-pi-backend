package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"log"
	"context"
	"crypto/tls"
	"time"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"

	"github.com/gorilla/mux"
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

func getUnits(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "return all units")
	// TODO serialize units dict, return that
}

func getUnit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "return unit with id")
	// TODO get id, version, and beanID from request, then publish update msg on mqtt channel with that data
}
func createUnit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create a unit")
	// TODO get id, version, and beanID from request, then publish update msg on mqtt channel with that data
}
func updateUnit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "update a unit")
	// TODO get id, version, and beanID from request, then publish update msg on mqtt channel with that data
}
func deleteUnit(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "delete a unit")
	// TODO get id, version, and beanID from request, then publish update msg on mqtt channel with that data
}

func makeHTTPServer() *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/api/units/", getUnits).Methods("GET")
	router.HandleFunc("/api/units/{id}", getUnit).Methods("GET")
	router.HandleFunc("/api/units", createUnit).Methods("POST")
	router.HandleFunc("/api/units/{id}", updateUnit).Methods("PUT")
	router.HandleFunc("/api/units/{id}", deleteUnit).Methods("DELETE")

	router.Handle("/lab/", http.StripPrefix("/lab/", http.FileServer(http.Dir("./static"))))
	//router.PathPrefix("/lab/{labID}").Handler(http.FileServer(http.Dir("./static")))

	return &http.Server{
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
        IdleTimeout:  120 * time.Second,
        Handler:      router,
    }
}

func main() {
	log.Println("Starting Backend")
	
	dict = make(map[string]Unit)

	var m *autocert.Manager
	var server *http.Server

	hostPolicy := func(ctx context.Context, host string) error {
		allowedHost := "saturten.com"
		if host == allowedHost {
			return nil
		}

		return fmt.Errorf("acme/autocert: only %s host is allowed", allowedHost)
	}

	certPath := "/etc/letsencrypt/live/saturten.com/"
	m = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: hostPolicy,
		Cache:      autocert.DirCache(certPath),
	}

	server = makeHTTPServer()
	server.Addr = ":443"
	server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
	server.TLSConfig.NextProtos = append(server.TLSConfig.NextProtos, acme.ALPNProto)

	log.Println("Starting server on ", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}
