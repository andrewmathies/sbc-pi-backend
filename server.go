package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"bytes"
	"io"
	"io/ioutil"
	"encoding/json"
	"os"

	"github.com/gorilla/mux"
)

var dict map[string]Unit

type Unit struct {
	Version string `json:"version"`
	BeanID 	string `json:"beanID"`
}

// Utility

func log(message string) {
	fmt.Println(message + "\n")
}

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
	log(string(requestDump))
}

func unitHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	dict = make(map[string]Unit)
	
	log("Starting server")

	router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/unit/{unit}", unitHandler)
	log.Fatal(http.ListenAndServeTLS(":443", "server.crt", "server.key", router))
}