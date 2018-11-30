package cola

import (
	"errors"
	"fmt"
	"testing"

	"github.com/briyanadityatama/goLoans/lms"
	"github.com/briyanadityatama/goLoans/lms/cola/domain"
)
import "github.com/stretchr/testify/assert"

const (
	ktpNumber = "3522582509010002"
	amount    = 10000000
	term      = 30
	birthDate = "1 December 1994"
	name      = "Doe"
)

var clientData = lms.ClientData{KTPNumber: ktpNumber, BirthDate: birthDate, Name: name}

func TestLmsRegisterClient(t *testing.T) {
	clientRepo := NewFakeClientRepo()
	cola := New(clientRepo)
	client, err := cola.RegisterClient(clientData)
	t.Run("should create a new client", func(t *testing.T) {
		assert.Equal(t, client.KTPNumber(), ktpNumber)
		assert.Equal(t, client.BirthDate(), birthDate)
		assert.Equal(t, client.Name(), name)
	})
	t.Run("should save a new client in repo", func(t *testing.T) {
		clientFoundInRepo, _, _ := clientRepo.ByKTPNumber(ktpNumber)
		assert.Equal(t, client, clientFoundInRepo)
	})
	t.Run("error should be nil", func(t *testing.T) {
		assert.Nil(t, err)
	})
}

func TestLmsRegisterClientTwice(t *testing.T) {
	clientRepo := NewFakeClientRepo()
	cola := New(clientRepo)
	cola.RegisterClient(clientData)
	client, err := cola.RegisterClient(clientData)
	t.Run("should return error", func(t *testing.T) {
		assert.NotNil(t, err)
		assert.Equal(t, "client_already_exists", err.Error())
		assert.Equal(t, err, lms.ErrClientAlreadyExists)
	})
	t.Run("returned client should be the one found in database", func(t *testing.T) {
		clientFoundInRepo, _, _ := clientRepo.ByKTPNumber(ktpNumber)
		assert.Equal(t, client, clientFoundInRepo)
	})
}

func TestLmsClientByPersonalNumber(t *testing.T) {
	clientRepo := NewFakeClientRepo()
	cola := New(clientRepo)
	client := domain.NewClient(birthDate, name, ktpNumber)
	clientRepo.Save(client)
	returnedClient, found, _ := cola.ClientByKTPNumber(ktpNumber)
	assert.True(t, found)
	assert.Equal(t, client, returnedClient)
}

func TestLmsApplyForLoan(t *testing.T) {
	clientRepo := NewFakeClientRepo()
	cola := New(clientRepo)
	clientRepo.Save(domain.NewClient(birthDate, name, ktpNumber))
	// when
	err := cola.ApplyForLoan(ktpNumber, amount, term)
	t.Run("should not return error", func(t *testing.T) {
		assert.Nil(t, err)
	})
	t.Run("should create a new loan", func(t *testing.T) {
		client, _, _ := cola.ClientByKTPNumber(ktpNumber)
		assert.True(t, client.HasActiveLoan())
	})
}

func TestLmsApplyForLoanWhenClientDoesNotExist(t *testing.T) {
	clientRepo := NewFakeClientRepo()
	cola := New(clientRepo)
	err := cola.ApplyForLoan(ktpNumber, amount, term)
	assert.Equal(t, err, lms.ErrClientDoesNotExist)
}

type SaveFailingClientRepo struct {
	ClientRepo
}

func (repo *SaveFailingClientRepo) Save(client domain.Client) error {
	return errors.New("database is down")
}

func TestLmsWhenRepoSavingIsFailing(t *testing.T) {
	failingClientRepo := &SaveFailingClientRepo{ClientRepo: NewFakeClientRepo()}
	cola := New(failingClientRepo)
	t.Run("RegisterClient", func(t *testing.T) {
		_, err := cola.RegisterClient(clientData)
		expectedErr := fmt.Sprintf("registering client %s %s with ktp number %s: database is down", birthDate, name, ktpNumber)
		assert.Equal(t, expectedErr, err.Error())
	})
}

type byKTPFailingClientRepo struct {
	ClientRepo
}

func (repo *byKTPFailingClientRepo) ByKTPNumber(ktpNumber string) (client domain.Client, found bool, err error) {
	return nil, false, errors.New("database is down again")
}

func TestLmsWhenRepoByPersonalNumberIsFailing(t *testing.T) {
	failingClientRepo := &byKTPFailingClientRepo{ClientRepo: NewFakeClientRepo()}
	cola := New(failingClientRepo)
	t.Run("RegisterClient", func(t *testing.T) {
		_, err := cola.RegisterClient(clientData)
		expectedErr := fmt.Sprintf("registering client %s %s with ktp number %s: database is down again", birthDate, name, ktpNumber)
		assert.Equal(t, expectedErr, err.Error())
	})
	t.Run("ClientByKTPNumber", func(t *testing.T) {
		_, _, err := cola.ClientByKTPNumber(ktpNumber)
		assert.Equal(t, "loading client by ktp number "+ktpNumber+": database is down again", err.Error())
	})
	t.Run("ApplyForLoan", func(t *testing.T) {
		err := cola.ApplyForLoan(ktpNumber, amount, term)
		expectedErr := fmt.Sprintf("client %s is applying for %d loan with term %d: database is down again", ktpNumber, amount, term)
		assert.Equal(t, expectedErr, err.Error())
	})
}
