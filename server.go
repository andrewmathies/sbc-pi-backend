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

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var dict map[string]Unit

type State int

const (
	Idle		State = 0
	Updating	State = 1
	Failed		State = 2
)

type Unit struct {
	Version string 	`json:"version"`
	BeanID 	string 	`json:"beanID"`
	Name	string	`json:"name"`
	State	State	`json:"state"`
}

// header is one of: Hello, StartUpdate, Complete, Fail
type Msg struct {
	ID		string	`json:"id"`
	Header	string	`json:"header"`
	Version	string	`json:"version"`
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
	dict[ksuid.New().String()] = Unit{Version: "2.3.4.5", BeanID: "12123434", Name: "ps960", State: 0}
	dict[ksuid.New().String()] = Unit{Version: "2.12.44.0", BeanID: "00009999", Name: "mangoooo", State: 1}
	dict[ksuid.New().String()] = Unit{Version: "1.9.8.7", BeanID: "98765432", Name: "PKD7000", State: 2}
	dict[ksuid.New().String()] = Unit{Version: "2.27", BeanID: "44553322", Name: "insert fake name here", State: 1}
}

// MQTT stuff

var client MQTT.Client
var qos int

func handleMsg(beanID string, msg Msg) {
	log.Println("Handling MQTT Message\nBeanID: ", beanID, "\nMsg: ", msg)

	switch msg.Header {
	case "Hello":
		// create new unit and stuff it in the dict
		unit := Unit{Version: msg.Version, BeanID: beanID, Name: "", State: Idle}
		dict[ksuid.New().String()] = unit
	case "StartUpdate":
		// backend should publish this, not recieve it
	case "Complete":
		// update status of unit and push that to frontend???
		id := msg.ID
		unit := dict[id]
		unit.State = Idle
		dict[id] = unit
	case "Fail":
		// update status of unit and push that to frontend???
		id := msg.ID
		unit := dict[id]
		unit.State = Failed
		dict[id] = unit
	default:
		log.Println("ERROR: unexpected MQTT message ", msg.Header)
	}
}

func publishMsg(beanID string, msg Msg) {
	log.Println("Publishing StartUpdate version on /unit/", beanID, "/")
	json, encodeErr := json.Marshal(msg)
	
	if encodeErr != nil {
		log.Println("ERROR: couldn't marshal msg for mqtt message ", encodeErr)
		return
	}

    token := client.Publish("/unit/" + beanID + "/", byte(qos), false, string(json))
    token.Wait()
}

// default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, mqttMsg MQTT.Message) {
	var msg Msg
	unpackErr := json.Unmarshal(mqttMsg.Payload(), &msg)
	if unpackErr != nil {
		log.Println("ERROR: couldn't unpack MQTT message ", unpackErr)
		return
	}

	topic := mqttMsg.Topic()
	topicParts := strings.Split(topic, "/")
	if len(topicParts) != 4 {
		log.Println("ERROR: badly formed MQTT topic ", topicParts)
		return
	}

	beanID := topicParts[2]
	handleMsg(beanID, msg)
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
	client = MQTT.NewClient(opts)
  	if token := client.Connect(); token.Wait() && token.Error() != nil {
    	panic(token.Error())
	}
	log.Println("Connected to MQTT broker")
	
	// subscribe to wildcard topic
	if token := client.Subscribe("/unit/+/", byte(qos), nil); token.Wait() && token.Error() != nil {
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
	oldUnit := dict[params["id"]]

	var temp Unit
	_ = json.NewDecoder(r.Body).Decode(&temp)

	if (temp.Version != "" && oldUnit.Version != temp.Version) {
		temp.State = Updating

		var msg Msg
		msg.ID = params["id"]
		msg.Header = "StartUpdate"
		msg.Version = temp.Version
		go publishMsg(oldUnit.BeanID, msg)
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

func makeHTTPServer() *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/api/units/", getUnits).Methods("GET")
	router.HandleFunc("/api/units/{id}", getUnit).Methods("GET")
	router.HandleFunc("/api/units", createUnit).Methods("POST")
	router.HandleFunc("/api/units/{id}", updateUnit).Methods("PUT")
	router.HandleFunc("/api/units/{id}", deleteUnit).Methods("DELETE")

	router.Handle("/", http.FileServer(http.Dir("views")))
	router.Handle("/lab", http.FileServer(http.Dir("views/lab")))
	router.Handle("/static", http.FileServer(http.Dir("static")))

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
	qos = 1
	go setupMQTT(tlsConfig)

	// build and run the https server
	server = makeHTTPServer()
	server.Addr = ":443"
	server.TLSConfig = tlsConfig
	server.TLSConfig.NextProtos = append(server.TLSConfig.NextProtos, acme.ALPNProto)

	log.Println("Starting server on ", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}