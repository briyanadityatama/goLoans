// Package rest decorates standard http.Request and http.ResponseWriter with methods useful for writing REST applications
package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// HandlerFunc wraps handler parameter into standard http.HandlerFunc
func HandlerFunc(handler func(*ResponseWriter, *Request)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := &Request{r}
		writer := &ResponseWriter{w}
		handler(writer, request)
	})
}

// Request is http.Request decorator adding useful methods for writing REST applications
type Request struct {
	*http.Request
}

// ReadJSONBody unmarshals JSON from HTTP body into output parameter
func (request *Request) ReadJSONBody(output interface{}) error {
	bytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return fmt.Errorf("Problem reading from http body when handling %s\n: %s", request.URL.Path, err.Error())
	}
	if err := json.Unmarshal(bytes, output); err != nil {
		return fmt.Errorf("JSON unmarshaling failed for http request: %s", err.Error())
	}
	return nil
}

// ResponseWriter is http.ResponseWriter decorator adding useful methods for writing REST applications
type ResponseWriter struct {
	http.ResponseWriter
}

// WriteJSON marshals input parameter into JSON and writes it to the connection
func (writer *ResponseWriter) WriteJSON(input interface{}) error {
	bytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("JSON marshalling failed for http response: %s", err.Error())
	}
	_, err = writer.Write(bytes)
	if err != nil {
		return fmt.Errorf("JSON marshalling failed for http response: %s", err.Error())
	}
	return nil
}

// WriteJSONError marshals err into JSON and writes it to the connection along with HTTP status code
func (writer *ResponseWriter) WriteJSONError(err error, statusCode int) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	jsonError := jsonError{err.Error(), err}
	e := writer.WriteJSON(jsonError)
	if e != nil {
		log.Printf("[WARN] Problem during marshaling error to JSON: %s", e.Error())
	}
}

type jsonError struct {
	Error  string      `json:"error"`
	Params interface{} `json:"params"`
}
