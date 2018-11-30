// Package rest exports lms functionality as Rest API
package rest

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/briyanadityatama/goLoans/lms"
	"github.com/briyanadityatama/goLoans/rest/rest"
)

// LoansServer is REST server providing lms functionality
type LoansServer struct {
	addr      string
	publicURL string
	lms       lms.Lms
	server    *http.Server
}

// NewLoansServer initialize LoansServer
func NewLoansServer(addr string, publicURL string, lms lms.Lms) *LoansServer {
	return &LoansServer{addr: addr, publicURL: publicURL, lms: lms}
}

// Start blocks current goroutine
func (server *LoansServer) Start() error {
	log.Printf("Starting server on http://%s/", server.addr)
	mux := http.NewServeMux()
	mux.Handle("/", rest.HandlerFunc(func(writer *rest.ResponseWriter, request *rest.Request) {
		writer.WriteHeader(404)
		fmt.Fprintln(writer, "Use /clients")
	}))
	mux.Handle("/clients", rest.HandlerFunc(func(writer *rest.ResponseWriter, request *rest.Request) {
		switch request.Method {
		case "GET":
			server.getClients(writer, request)
		case "POST":
			server.postClients(writer, request)
		case "OPTIONS":
			server.optionsClients(writer, request)
		}
	}))
	mux.Handle("/clients/", rest.HandlerFunc(func(writer *rest.ResponseWriter, request *rest.Request) {
		switch request.Method {
		case "GET":
			server.getClient(writer, request)
		}
	}))
	server.server = &http.Server{Addr: server.addr, Handler: mux}
	return server.server.ListenAndServe()
}

// Stop can be executed from a different goroutine to stop the server
func (server *LoansServer) Stop() error {
	return server.server.Shutdown(nil)
}

func (server *LoansServer) getClients(w *rest.ResponseWriter, r *rest.Request) {
	w.WriteHeader(405)
	fmt.Fprintln(w, "POST /clients to register a new client")
}

func (server *LoansServer) postClients(writer *rest.ResponseWriter, request *rest.Request) {
	var clientData lms.ClientData
	err := request.ReadJSONBody(&clientData)
	if err != nil {
		writer.WriteHeader(400)
		fmt.Fprintln(writer, err.Error())
		return
	}
	client, err := server.lms.RegisterClient(clientData)
	if err != nil {
		writer.WriteJSONError(err, 400)
		return
	}
	writer.Header().Add("Location", server.publicURL+"/clients/"+client.KTPNumber())
	writer.WriteHeader(201)
}

func (server *LoansServer) getClient(writer *rest.ResponseWriter, request *rest.Request) {
	ktpNumber := request.URL.Path[len("/clients/"):]
	client, found, err := server.lms.ClientByKTPNumber(ktpNumber)
	if err != nil {
		errorDto := fmt.Sprintf("problem getting client with ktpNumber %s: %s", ktpNumber, err.Error())
		serverError := technicalError{errors.New("server_error"), errorDto}
		writer.WriteJSONError(serverError, 500)
		return
	}
	if !found {
		writer.WriteJSONError(lms.ErrClientDoesNotExist, 404)
		return
	}
	writer.WriteHeader(200)
	selfLink := fmt.Sprintf("%s/clients/%s/goLoans", server.publicURL, client.KTPNumber())
	response := getClientResponse{
		KTPNumber: client.KTPNumber(),
		BirthDate: client.BirthDate(),
		Name:      client.Name(),
		Loans:     goLoans{[]link{{"self", selfLink}}}}
	err = writer.WriteJSON(response)
	if err != nil {
		log.Printf("[WARN] problem getting client with ktpNumber %s: %s", ktpNumber, err.Error())
	}
}

// getClientResponse DTO for JSON marshaling
type getClientResponse struct {
	KTPNumber string  `json:"ktpNumber"`
	BirthDate string  `json:"birthDate"`
	Name      string  `json:"name"`
	Loans     goLoans `json:"goLoans"`
}

// goLoans DTO for JSON marshaling
type goLoans struct {
	Links []link `json:"links"`
}

// link DTO for JSON marshaling
type link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

// technicalError DTO for JSON marshaling
type technicalError struct {
	error
	TechnicalError string
}

func (server *LoansServer) optionsClients(w *rest.ResponseWriter,
	r *rest.Request) {
	w.Header().Add("Allow", "OPTIONS, POST")
	w.WriteHeader(200)
}
