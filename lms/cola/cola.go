package cola

import (
	"fmt"

	"github.com/briyanadityatama/goLoans/lms"
	"github.com/briyanadityatama/goLoans/lms/cola/domain"
)

// ClientRepo is used internally by lms package for loading/storing client information. This makes lms package
// independent of database technology.
type ClientRepo interface {
	ByKTPNumber(ktpNumber string) (client domain.Client, found bool, err error)
	Save(client domain.Client) error
}

type cola struct {
	ClientRepo ClientRepo
}

// New returns a new instance of Lms
func New(repo ClientRepo) lms.Lms {
	return &cola{ClientRepo: repo}
}

func (cola *cola) RegisterClient(clientData lms.ClientData) (lms.Client, error) {
	gender := clientData.Gender
	birthDate := clientData.BirthDate
	name := clientData.Name
	ktpNumber := clientData.KTPNumber
	client, found, err := cola.ClientRepo.ByKTPNumber(ktpNumber)
	if err != nil {
		return nil, fmt.Errorf("registering client %s %s with personal number %s: %v", gender, birthDate, name, ktpNumber, err)
	}
	if found {
		return client, lms.ErrClientAlreadyExists
	}
	client = domain.NewClient(gender, birthDate, name, ktpNumber)
	err = cola.ClientRepo.Save(client)
	if err != nil {
		return nil, fmt.Errorf("registering client %s %s with personal number %s: %v", gender, birthDate, name, ktpNumber, err)
	}
	return client, nil
}

func (cola *cola) ClientByKTPNumber(ktpNumber string) (client lms.Client, found bool, err error) {
	client, found, err = cola.ClientRepo.ByKTPNumber(ktpNumber)
	if err != nil {
		err = fmt.Errorf("loading client by personal number %s: %v", ktpNumber, err)
	}
	return
}
func (cola *cola) ApplyForLoan(ktpNumber string, amount uint, term uint) error {
	client, found, err := cola.ClientRepo.ByKTPNumber(ktpNumber)
	if err != nil {
		return fmt.Errorf("client %s is applying for %d loan with term %d: %v", ktpNumber, amount, term, err)
	}
	if !found {
		return lms.ErrClientDoesNotExist
	}
	applicationError := client.ApplyForLoan(amount, domain.Term(term))
	return applicationError
}
