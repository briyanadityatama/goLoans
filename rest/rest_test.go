package rest

import (
	"errors"
	"strings"
	"testing"

	"github.com/briyanadityatama/goLoans/lms"
	"github.com/briyanadityatama/goLoans/testing/http"
	"github.com/stretchr/testify/assert"
)

const (
	ktpNumber = "3522582509010002"
	birthDate = "1 December 1994"
	name      = "Doe"
)

var clientData = lms.ClientData{KTPNumber: ktpNumber, BirthDate: birthDate, Name: name}

func TestGetSlash(t *testing.T) {
	server := newServer(lms.NewFakeLms())
	go server.Start()
	defer server.Stop()
	response, status := http.Get("/")
	assert.Equal(t, 404, status)
	assert.Equal(t, "Use /clients", response)
}

func newServer(lmsImplementation lms.Lms) *LoansServer {
	return NewLoansServer(http.Address, "http://"+http.Address, lmsImplementation)
}

func TestClients(t *testing.T) {
	fakeLms := lms.NewFakeLms()
	server := newServer(fakeLms)
	go server.Start()
	defer server.Stop()
	t.Run("GET /clients", func(t *testing.T) {
		response, status := http.Get("/clients")
		assert.Equal(t, 405, status)
		assert.Equal(t, "POST /clients to register a new client", response)
	})
	t.Run("OPTIONS /clients", func(t *testing.T) {
		allowHeader, status := http.Options("/clients")
		assert.Equal(t, 200, status)
		assert.Equal(t, "OPTIONS, POST", allowHeader)
	})
	t.Run("POST /clients", func(t *testing.T) {
		_, status, headers := http.Post("/clients",
			`{
			"ktpNumber": "`+ktpNumber+`",
			"birthDate": "`+birthDate+`",
			"name": "`+name+`"
					}`)
		t.Run("Should return 201 with Location header", func(t *testing.T) {
			assert.Equal(t, 201, status)
			assert.Equal(t, server.publicURL+"/clients/"+ktpNumber, headers.Get("Location"))
		})
		t.Run("Should register a new client in lms", func(t *testing.T) {
			client, clientFound, _ := fakeLms.ClientByKTPNumber(ktpNumber)
			assert.True(t, clientFound)
			if clientFound {
				assert.Equal(t, client.KTPNumber(), ktpNumber)
				assert.Equal(t, client.BirthDate(), birthDate)
				assert.Equal(t, client.Name(), name)
			}
		})
	})
	t.Run("GET /clients/{ktpNumber}", func(t *testing.T) {
		fakeLms.RegisterClient(clientData)
		response, status := http.Get("/clients/" + ktpNumber)
		assert.Equal(t, 200, status)
		expectedResponse := map[string]interface{}{
			"ktpNumber": ktpNumber,
			"birthDate": birthDate,
			"name":      name,
			"goLoans": map[string]interface{}{
				"links": []interface{}{
					map[string]interface{}{
						"rel":  "self",
						"href": server.publicURL + "/clients/" + ktpNumber + "/goLoans",
					},
				},
			},
		}
		assert.Equal(t, expectedResponse, http.Unmarshal(response))
	})
	t.Run("GET /clients/{unexistingKTPNumber}", func(t *testing.T) {
		response, status := http.Get("/clients/1")
		assert.Equal(t, 404, status)
		expectedResponse := map[string]interface{}{
			"error":  "client_does_not_exist",
			"params": map[string]interface{}{},
		}
		assert.Equal(t, expectedResponse, http.Unmarshal(response))
	})
}

type LmsFailingOnRegistration struct {
	lms.Lms
}

func (*LmsFailingOnRegistration) RegisterClient(clientData lms.ClientData) (lms.Client, error) {
	return nil, lms.ErrClientAlreadyExists
}

func TestPostClientWhenClientAlreadyExists(t *testing.T) {
	// given
	server := newServer(&LmsFailingOnRegistration{lms.NewFakeLms()})
	go server.Start()
	defer server.Stop()
	// when
	response, status, headers := http.Post("/clients",
		`{
			"ktpNumber": "3522582509010001",
			"birthDate": "2 December 1994",
			"name": "Bar"
					}`)
	// then
	assert.Equal(t, 400, status)
	expectedResponse := map[string]interface{}{
		"error":  "client_already_exists",
		"params": map[string]interface{}{},
	}
	assert.Equal(t, expectedResponse, http.Unmarshal(response))
	assert.Equal(t, "application/json", headers.Get("Content-Type"))
}

func TestPostClientWithIncorrectJson(t *testing.T) {
	server := newServer(lms.NewFakeLms())
	go server.Start()
	defer server.Stop()
	bodies := []string{"", "[]", " ", "\n", "{"}
	for _, body := range bodies {
		response, status, _ := http.Post("/clients", body)
		assert.Equal(t, 400, status)
		assert.True(t, strings.HasPrefix(response, "JSON unmarshaling failed"))
	}
}

type LmsFailingOnClientByKTPNumber struct {
	lms.Lms
}

func (*LmsFailingOnClientByKTPNumber) ClientByKTPNumber(ktpNumber string) (client lms.Client, found bool, error error) {
	return nil, false, errors.New("clientByKTPNumber failed")
}

func TestGetClientWhenLmsIsFailing(t *testing.T) {
	// given
	fakeLms := lms.NewFakeLms()
	fakeLms.RegisterClient(clientData)
	server := newServer(&LmsFailingOnClientByKTPNumber{fakeLms})
	go server.Start()
	defer server.Stop()
	// when
	response, status := http.Get("/clients/" + ktpNumber)
	// then
	assert.Equal(t, 500, status)
	expectedResponse := map[string]interface{}{
		"error": "server_error",
		"params": map[string]interface{}{
			"TechnicalError": "problem getting client with ktpNumber " + ktpNumber + ": clientByKTPNumber failed",
		},
	}
	assert.Equal(t, expectedResponse, http.Unmarshal(response))
}
