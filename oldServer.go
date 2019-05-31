package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"bytes"
	"io"
	"io/ioutil"
	"encoding/json"
	"encoding/csv"
	"os"

	"github.com/gorilla/mux"
)

// Globals----------------------------------------------------------------

const csvPath = "/go/sbc-pi-backend/static/versions.csv"

var dict map[string]Unit

// Message Structs--------------------------------------------------------

type AddMsg struct {
	Header string `json:"header"`
	ID string `json:"id"`
	BeanID string `json:"beanID"`
}

type RemoveMsg struct {
	Header string `json:"header"`
	ID string `json:"id"`
}

type UpdateMsg struct {
	Header string `json:"header"`
	Version string `json:"version"`
	ID string `json:"id"`
}

type Unit struct {
	Version string `json:"version"`
	BeanID 	string `json:"beanID"`
}

// Enum-------------------------------------------------------------------

type Msg int

const (
	Add	Msg = 0
	Remove	Msg = 1
	Update	Msg = 2
)

// Utility----------------------------------------------------------------

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

func routeToPi(r io.Reader) {
	/*
	if piConn != nil {
		body, err := ioutil.ReadAll(r)
		checkError("routeToPi, reading body", err)
	
		readBuf := bytes.NewBuffer(body)
		piConn.Write(readBuf.Bytes())
	}
	*/
}

func getUnits(id string, version string, operation Msg) [][]string {
	var units [][]string

	file, err := os.OpenFile(csvPath, os.O_RDONLY, 0666)
	checkError("cannot open file", err)
	reader := csv.NewReader(file)

	for {
		record, readErr := reader.Read()
		
		if readErr == io.EOF {
			break
		}
		checkError("reading csv", readErr)

		if operation == Update && record[0] == id {
			record[1] = version
		} else if operation == Remove && record[0] == id {
			continue
		}
		
		units = append(units, record)
	}
	
	deleteErr := os.Remove(csvPath)
	checkError("deleting csv file at: " + csvPath, deleteErr)

	return units
}

func saveCSV(units [][]string) {
	newFile, createErr := os.Create(csvPath)
	checkError("creating csv file", createErr)
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	defer writer.Flush()

	for _, line := range units {
		writeErr := writer.Write(line)
		checkError("writing units to csv file", writeErr)		
	}

	log("units written to disk")
}

// HTTP Endpoints---------------------------------------------------------

func unitHandler(w http.ResponseWriter, r *http.Request) {

}

func updateVersion(w http.ResponseWriter, r *http.Request) {
	log("update version endpoint hit")	
	formatRequest(r)

	var copyBuf bytes.Buffer
	tee := io.TeeReader(r.Body, &copyBuf)
	
	body, parseError := ioutil.ReadAll(tee)
	checkError("grabbing raw data from response body", parseError)

	msg := UpdateMsg{}
	marshalErr := json.Unmarshal(body, &msg)
	checkError("parsing json from raw data", marshalErr)
	
	routeToPi(&copyBuf)
	
	units := getUnits(msg.ID, msg.Version, Update)
	saveCSV(units)
}


func addUnit(w http.ResponseWriter, r *http.Request) {
	log("add unit endpoint hit")
	formatRequest(r)

	var copyBuf bytes.Buffer
	tee := io.TeeReader(r.Body, &copyBuf)

	body, parseError := ioutil.ReadAll(tee)
	checkError("grabbing raw data from response body", parseError)

	msg := AddMsg{}
	marshalErr := json.Unmarshal(body, &msg)
	checkError("parsing json from raw data", marshalErr)
	
	unit = Unit{msg.BeanID, ""}

	routeToPi(&copyBuf)

	file, err := os.OpenFile(csvPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	checkError("opening csv file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{ msg.ID, "", msg.BeanID })
	log("add unit change written to disk")
}

func removeUnit(w http.ResponseWriter, r *http.Request) {
	log("remove unit endpoint hit")
	formatRequest(r)

	var copyBuf bytes.Buffer
	tee := io.TeeReader(r.Body, &copyBuf)

	body, parseError := ioutil.ReadAll(tee)
	checkError("grabbing raw data from response body", parseError)

	msg := RemoveMsg{}
	marshalErr := json.Unmarshal(body, &msg)
	checkError("parsing json from raw data", marshalErr)
	
	routeToPi(&copyBuf)

	units := getUnits(msg.ID, "", Remove)
	saveCSV(units)
}

// Server-----------------------------------------------------------------

func setupRoutes() {
	http.Handle("/", http.FileServer(http.Dir("static/")))
	http.HandleFunc("/unit", unitHandler)
	
	err := http.ListenAndServe(":80", nil)
	checkError("http server crashed", err)
}

func main() {
	dict = make(map[string]Unit)
	log("Starting server")
	setupRoutes()
}
