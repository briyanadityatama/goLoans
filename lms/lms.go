// Package lms (Loans Management System) is an entry API for all application logic. It provides methods for every use case.
// In DDD terminology its an application layer. It interacts with database using repository abstractions to load and
// persist client with their loans to disk. This makes the package independent of any technology used for storing data.
// The package should be used by servers (such as ReST HTTP server) or graphical user interfaces.
package lms

import (
	"errors"
)

// Lms provides methods for all use cases in the system
type Lms interface {
	RegisterClient(clientData ClientData) (Client, error)
	ClientByKTPNumber(ktpNumber string) (client Client, found bool, error error)
	ApplyForLoan(ktpNumber string, amount uint, term uint) (error error)
}

// Client is someone who wants to take a loan
type Client interface {
	Gender() string
	KTPNumber() string
	BirthDate() string
	Name() string
	HasActiveLoan() bool
}

// ClientData stores personal information about client and is used as data transfer object DTO
type ClientData struct {
	Gender    string
	KTPNumber string
	BirthDate string
	Name      string
}

// ErrClientAlreadyExists is an error returned when Client already exists
var ErrClientAlreadyExists = errors.New("client_already_exists")

// ErrClientDoesNotExist is an error return when Client does not exist
var ErrClientDoesNotExist = errors.New("client_does_not_exist")
