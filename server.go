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
	"os"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	MQTT "github.com/eclipse/paho.mqtt.golang"
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

// MQTT stuff

// default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Println("TOPIC: %s\n", msg.Topic())
	log.Println("MSG: %s\n", msg.Payload())
}

func setupMQTT(tlsConfig *tls.Config) {
	// opts contains broker address and other config info
	opts := MQTT.NewClientOptions().AddBroker("tls://saturten.com:8883")
  	opts.SetClientID("go-simple")
	opts.SetDefaultPublishHandler(f)
	opts.SetTLSConfig(tlsConfig)
	opts.SetUsername("andrew")
	opts.SetPassword("1plus2is3")
	
	// initiate connection with broker
	c := MQTT.NewClient(opts)
  	if token := c.Connect(); token.Wait() && token.Error() != nil {
    	panic(token.Error())
	}
	log.Println("Connected to MQTT broker")
	
	// subscribe to wildcard topic
	if token := c.Subscribe("/unit/+/", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	log.Println("Subscribed to /unit/+/")
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
	delete(dict, params["id"])
	var unit Unit
	_ = json.NewDecoder(r.Body).Decode(&unit)
	dict[params["id"]] = unit
	json.NewEncoder(w).Encode(unit)
}

func deleteUnit(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE - deleteUnit hit")
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

	return &http.Server{
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 5 * time.Second,
        IdleTimeout:  120 * time.Second,
        Handler:      router,
    }
}

func main() {
	log.Println("Starting Backend")
	
	// data TODO: implement db
	dict = make(map[string]Unit)
	fakeData()

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

	tlsConfig := &tls.Config{GetCertificate: m.GetCertificate}

	// mqtt client
	go setupMQTT(tlsConfig)

	// build and run the https server
	server = makeHTTPServer()
	server.Addr = ":443"
	server.TLSConfig = tlsConfig
	server.TLSConfig.NextProtos = append(server.TLSConfig.NextProtos, acme.ALPNProto)

	log.Println("Starting server on ", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}
