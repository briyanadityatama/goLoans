// Package http provides functions for writing HTTP server tests in a concise manner
package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

// Address with available port which can be used for starting a server
var Address = availableAddr()

func availableAddr() (addr string) {
	server, err := net.Listen("tcp", ":0")
	if err != nil {
		panic("No ports are available!")
	}
	defer server.Close()
	hostString := server.Addr().String()
	return hostString
}

// Get runs HTTP GET method
func Get(path string) (responseBody string, status int) {
	response := do("GET", path, nil)
	responseBody = readResponseBody(response)
	status = response.StatusCode
	return
}

func readResponseBody(response *http.Response) string {
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Panic("Problem reading body", err)
	}
	return string(bytes)
}

// Options runs HTTP OPTIONS method
func Options(path string) (allowHeader string, status int) {
	response := do("OPTIONS", path, nil)
	allowHeader = response.Header.Get("Allow")
	status = response.StatusCode
	return
}

// Post runs HTTP POST method
func Post(path string, body string) (responseBody string, status int, header http.Header) {
	response := do("POST", path, strings.NewReader(body))
	responseBody = readResponseBody(response)
	status = response.StatusCode
	header = response.Header
	return
}

func do(method, path string, body io.Reader) (response *http.Response) {
	url := "http://" + Address + path
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Panicf("http request creation failed %s %s: %s", method, path, err)
	}
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		log.Panicf("http request failed for %s %s: %s", method, path, err)
	}
	return response
}

// Unmarshal the string and return map
func Unmarshal(response string) map[string]interface{} {
	returnedError := make(map[string]interface{})
	json.Unmarshal([]byte(response), &returnedError)
	return returnedError
}
