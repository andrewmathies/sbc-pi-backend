package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"html"
	"log"

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

func unitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. unit handler here")
	// TODO get id, version, and beanID from request, then publish update msg on mqtt channel with that data
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello. list handler here")
	// TODO convert dict to json, send that as response
}

func main() {
	dict = make(map[string]Unit)
	
	log.Println("Starting server")

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/unit/{beanID}", unitHandler)
	router.HandleFunc("/units/{labID}", listHandler)
	router.PathPrefix("/lab/{labID}").Handler(http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", router))
}