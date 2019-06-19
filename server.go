package main

import (
	"net/http"
	"net/http/httputil"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"context"
	"time"
	"os"
	"strings"

	"bg"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

// Globals

type Unit bg.Unit
type Msg bg.Msg
type State int

const (
	Idle		State = 0
	Updating	State = 1
	Failed		State = 2
)

var dict map[string]Unit

// Utility Functions

func checkErr(message string, err error) {
	if err != nil {
		log.Println("ERROR: " + message)
		log.Println(err)
	}
}

func formatRequest(r *http.Request) {
	requestDump, formatErr := httputil.DumpRequest(r, true)
	checkErr("couldnt format http request", formatErr)
	log.Println(string(requestDump))
}

func fakeData() {
	dict = make(map[string]Unit)

	dict[ksuid.New().String()] = Unit{Version: "2.3.4.5", BeanID: "12123434", Name: "ps960", State: 0}
	dict[ksuid.New().String()] = Unit{Version: "2.12.44.0", BeanID: "00009999", Name: "mangoooo", State: 1}
	dict[ksuid.New().String()] = Unit{Version: "1.9.8.7", BeanID: "98765432", Name: "PKD7000", State: 2}
	dict[ksuid.New().String()] = Unit{Version: "2.27", BeanID: "44553322", Name: "insert fake name here", State: 1}
}

// HTTP Handlers

func getUnits(w http.ResponseWriter, r *http.Request) {
	log.Println("GET - getUnits hit")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dict)
}

func getUnit(w http.ResponseWriter, r *http.Request) {
	log.Println("GET - getUnit hit")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	unit := dict[params["id"]]
	json.NewEncoder(w).Encode(unit)
}

func createUnit(w http.ResponseWriter, r *http.Request) {
	log.Println("POST - createUnit hit")
	w.Header().Set("Content-Type", "application/json")
	var unit Unit
	_ = json.NewDecoder(r.Body).Decode(&unit)
	id := ksuid.New().String()
	dict[id] = unit
	json.NewEncoder(w).Encode(unit)
}

func updateUnit(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT - updateUnit hit")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	oldUnit := dict[params["id"]]

	var temp Unit
	decodeErr := json.NewDecoder(r.Body).Decode(&temp)

	if decodeErr != nil {
		log.Println("PUT - failed decoding request")
		log.Println(decodeErr)
		json.NewEncoder(w).Encode(oldUnit)
		return
	}

	if (temp.Version != "" && oldUnit.Version != temp.Version) {
		temp.State = Updating

		var msg Msg
		msg.ID = params["id"]
		msg.Header = "StartUpdate"
		msg.Version = temp.Version
		go bg.publishMsg(oldUnit.BeanID, msg)
	}

	dict[params["id"]] = temp
	
	json.NewEncoder(w).Encode(temp)
}

func deleteUnit(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE - deleteUnit hit")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	delete(dict, params["id"])
	json.NewEncoder(w).Encode(dict)
}

func getTLSConfig() *tls.Config {
	// this section makes sure we have a valid cert
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

	tlsConfig := &tls.Config{ GetCertificate: m.GetCertificate }
	tlsConfig.NextProtos = append(tlsConfig.NextProtos, acme.ALPNProto)

	return tlsConfig
}

func httpServer(tlsConfig *tls.Config) {
	router := mux.NewRouter()

	// REST API
	router.HandleFunc("/api/units/", getUnits).Methods("GET")
	router.HandleFunc("/api/units/{id}", getUnit).Methods("GET")
	router.HandleFunc("/api/units", createUnit).Methods("POST")
	router.HandleFunc("/api/units/{id}", updateUnit).Methods("PUT")
	router.HandleFunc("/api/units/{id}", deleteUnit).Methods("DELETE")

	// FRONTEND
	router.PathPrefix("/lab/").Handler(http.StripPrefix("/lab/", http.FileServer(http.Dir("lab/"))))

	server := &http.Server{
    	ReadTimeout:	5 * time.Second,
        WriteTimeout:	5 * time.Second,
        IdleTimeout:	120 * time.Second,
		Handler:		router,
		Addr: 			":443"
		TLSConfig: 		tlsConfig
	}

	log.Println("Starting server on ", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func main() {
	log.Println("Starting Backend")
	
	// TODO: implement db
	fakeData()

	tlsConfig := getTLSConfig()

	go bg.initMQTT(tlsConfig)
	httpServer(tlsConfig)
}