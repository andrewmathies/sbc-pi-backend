package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"log"
	"context"
	"crypto/tls"
	"time"
	"encoding/json"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
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


func fakeData() {
	dict[ksuid.New().String()] = Unit{Version: "2.3.4.5", BeanID: "12123434"}
	dict[ksuid.New().String()] = Unit{Version: "2.12.44.0", BeanID: "00009999"}
	dict[ksuid.New().String()] = Unit{Version: "1.9.8.7", BeanID: "98765432"}
	dict[ksuid.New().String()] = Unit{Version: "2.27", BeanID: "44553322"}
}

// HTTP Handlers

func getUnits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dict)
}

func getUnit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	unit := dict[params["id"]]
	json.NewEncoder(w).Encode(unit)
}

func createUnit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var unit Unit
	_ = json.NewDecoder(r.Body).Decode(&unit)
	id := ksuid.New().String()
	dict[id] = unit
	json.NewEncoder(w).Encode(unit)
}

func updateUnit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	delete(dict, params["id"])
	_ = json.NewDecoder(r.Body).Decode(&unit)
	dict[params["id"]] = unit
	json.NewEncoder(w).Encode(unit)
}

func deleteUnit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	delete(dict, params["id"])
	json.NewEncoder(w).Encode(dict)
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
	fakeData()

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
