package main

import (
	"net/http"
	"net/http/httputil"
	"crypto/tls"
	"encoding/json"
	"hash/fnv"
	"fmt"
	"log"
	"context"
	"time"
	"os"
	"strings"
	"strconv"

	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"github.com/gorilla/mux"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var dict map[uint64]Unit
var versions map[uint64]string

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
	ID		uint64	`json:"id"`
	Header	string	`json:"header"`
	Version	string	`json:"version"`
}

// Utility

func checkErr(message string, err error) {
	if err != nil {
		fmt.Println("ERROR: " + message)
		fmt.Println(err)
		fmt.Println("")
	}
}

func formatRequest(r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	checkErr("couldnt format http request", err)
	log.Println(string(requestDump))
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func convert(s string) uint64 {	
	num, convErr := strconv.ParseUint(s, 10, 64)
	checkErr("converting string to uint64", convErr)
	return num
}

func fakeData() {
	dict = make(map[uint64]Unit)
	versions = make(map[uint64]string)
/*
	dict[hash("12123434")] = Unit{Version: "2.3.4.5", BeanID: "12123434", Name: "ps960", State: 0}
	dict[hash("00009999")] = Unit{Version: "2.12.44.0", BeanID: "00009999", Name: "mangoooo", State: 1}
	dict[hash("98765432")] = Unit{Version: "1.9.8.7", BeanID: "98765432", Name: "PKD7000", State: 2}
	dict[hash("44553322")] = Unit{Version: "2.27", BeanID: "44553322", Name: "insert fake name here", State: 1}
*/
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
		key := hash(beanID)
		if val, ok := dict[key]; ok {
			log.Println("key for bean ID: " + val.BeanID + " already exists")
		} else {
			dict[key] = unit
		}
	case "StartUpdate":
		// the backend published this, so do nothing
	case "Success":
		// update status of unit and push that to frontend???
		key := hash(beanID)
		unit := dict[key]
		unit.State = Idle
		dict[id] = unit
	case "Fail":
		// update status of unit and push that to frontend???
		key := hash(beanID)
		unit := dict[key]
		unit.State = Failed
		dict[id] = unit
	default:
		log.Println("ERROR: unexpected MQTT message ", msg.Header)
	}
}

func publishMsg(beanID string, msg Msg) {
	log.Println("Publishing StartUpdate version on unit/", beanID, "/")
	json, encodeErr := json.Marshal(msg)
	
	if encodeErr != nil {
		log.Println("ERROR: couldn't marshal msg for mqtt message ", encodeErr)
		return
	}

    token := client.Publish("unit/" + beanID + "/", byte(qos), false, string(json))
    token.Wait()
}

// default message handler
var f MQTT.MessageHandler = func(client MQTT.Client, mqttMsg MQTT.Message) {
	topic := mqttMsg.Topic()
	topicParts := strings.Split(topic, "/")

	log.Println("recieved msg on topic: " + topicParts[0])
	
	if topicParts[0] == "version" {
		key := hash(topicParts[1])
		versions[key] = topicParts[1]
		log.Println("adding " + topicParts[1] + " to versions map")
		return
	}

	if len(topicParts) != 3 {
		log.Println("ERROR: badly formed MQTT topic ", topicParts)
		return
	}
	
	var msg Msg
	unpackErr := json.Unmarshal(mqttMsg.Payload(), &msg)
	if unpackErr != nil {
		log.Println("ERROR: couldn't unpack MQTT message ", unpackErr)
		return
	}

	beanID := topicParts[1]
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

	qos = 1
	
	// initiate connection with broker
	client = MQTT.NewClient(opts)
  	if token := client.Connect(); token.Wait() && token.Error() != nil {
    	panic(token.Error())
	}
	log.Println("Connected to MQTT broker")
	
	// subscribe to unit wildcard topic
	if token := client.Subscribe("unit/+/", byte(qos), nil); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}
	log.Println("Subscribed to unit/+/")

	// subscribe to version wildcard topic
	if token := client.Subscribe("version/+/", byte(qos), nil); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}
	log.Println("Subscribed to version/+/")
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
	unit := dict[convert(params["id"])]
	json.NewEncoder(w).Encode(unit)
}

func createUnit(w http.ResponseWriter, r *http.Request) {
	log.Println("POST - createUnit hit")
	w.Header().Set("Content-Type", "application/json")
	var unit Unit
	_ = json.NewDecoder(r.Body).Decode(&unit)
	id := hash("test")
	dict[id] = unit
	json.NewEncoder(w).Encode(unit)
}

func updateUnit(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT - updateUnit hit")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	key := convert(params["id"])
	oldUnit := dict[key]

	var reqData Unit
	decodeErr := json.NewDecoder(r.Body).Decode(&reqData)

	if decodeErr != nil {
		checkErr("PUT - failed decoding request", decodeErr)
		json.NewEncoder(w).Encode(oldUnit)
		return
	}

	if (reqData.Version != "" && oldUnit.Version != reqData.Version) {
		reqData.State = Updating

		var msg Msg
		msg.ID = key
		msg.Header = "StartUpdate"
		msg.Version = reqData.Version
		go publishMsg(oldUnit.BeanID, msg)
	}

	dict[key] = reqData
	
	json.NewEncoder(w).Encode(reqData)
}

func deleteUnit(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE - deleteUnit hit")
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	delete(dict, convert(params["id"]))
	json.NewEncoder(w).Encode(dict)
}

func getVersions(w http.ResponseWriter, r *http.Request) {
	log.Println("GET - get versions endpoint hit")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(versions)
}

func makeHTTPServer(tlsConfig *tls.Config) *http.Server {
	router := mux.NewRouter()
	router.HandleFunc("/api/units/", getUnits).Methods("GET")
	router.HandleFunc("/api/units/{id}", getUnit).Methods("GET")
	router.HandleFunc("/api/units", createUnit).Methods("POST")
	router.HandleFunc("/api/units/{id}", updateUnit).Methods("PUT")
	router.HandleFunc("/api/units/{id}", deleteUnit).Methods("DELETE")
	router.HandleFunc("/api/versions/", getVersions).Methods("GET")

	router.PathPrefix("/lab/").Handler(http.StripPrefix("/lab/", http.FileServer(http.Dir("lab/"))))

	return &http.Server{
        ReadTimeout:  	5 * time.Second,
        WriteTimeout: 	5 * time.Second,
        IdleTimeout:  	120 * time.Second,
		Handler:      	router,
		Addr:			":443",
		TLSConfig:		tlsConfig,
    }
}

func getTlsConfig() *tls.Config {
	// this section makes sure we have a valid cert
	var m *autocert.Manager

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

	config := &tls.Config{GetCertificate: m.GetCertificate}
	config.NextProtos = append(config.NextProtos, acme.ALPNProto)

	return config
}

func main() {
	log.Println("Starting Backend")
	
	// data TODO: implement db
	fakeData()

	tlsConfig := getTlsConfig()

	// mqtt client
	go setupMQTT(tlsConfig)

	// build and run the https server
	server := makeHTTPServer(tlsConfig)

	log.Println("Starting server on ", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}
