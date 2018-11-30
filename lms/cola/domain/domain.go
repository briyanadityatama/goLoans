// Package domain provides core business logic which is independent of any other systems and repositories
package domain

import "errors"

const maximumAmountForFirstLoan = 50000000

// Client can only have one active loan
type Client interface {
	KTPNumber() string
	BirthDate() string
	Name() string
	Gender() string
	ApplyForLoan(amount uint, term Term) (err error)
	HasActiveLoan() bool
	ActiveLoan() Loan
	Repay(amount uint) (err error)
}

// NewClient returns Client instance
func NewClient(gender, birthDate, name, ktpNumber string) Client {
	return &paydayLoanClient{gender: gender, birthDate: birthDate, name: name, ktpNumber: ktpNumber}
}

// Loan should be repaid in a given term or something bad will happen
type Loan interface {
	Amount() uint
	Term() Term
	Remaining() uint
}

type paydayLoan struct {
	amount    uint
	term      Term
	remaining uint
}

func (loan *paydayLoan) Remaining() uint {
	return loan.remaining
}

type paydayLoanClient struct {
	ktpNumber string
	birthDate string
	name      string
	gender    string
	loan      *paydayLoan
}

func (client *paydayLoanClient) ActiveLoan() Loan {
	return client.loan
}

func (client *paydayLoanClient) KTPNumber() string {
	return client.ktpNumber
}

func (client *paydayLoanClient) BirthDate() string {
	return client.birthDate
}

func (client *paydayLoanClient) Name() string {
	return client.name
}

func (client *paydayLoanClient) Gender() string {
	return client.gender
}

func (client *paydayLoanClient) HasActiveLoan() bool {
	return client.loan != nil
}

func (client *paydayLoanClient) ApplyForLoan(amount uint, term Term) error {
	if client.HasActiveLoan() {
		return ErrClientAlreadyHasLoan
	}
	if amount > maximumAmountForFirstLoan {
		return ErrAmountTooHigh
	}
	loan := &paydayLoan{amount: amount, term: term, remaining: amount}
	client.loan = loan
	return nil
}

func (client *paydayLoanClient) Repay(amount uint) (err error) {
	repaymentError := client.loan.repay(amount)
	if repaymentError != nil {
		return repaymentError
	}
	if client.loan.Remaining() == 0 {
		client.loan = nil
	}
	return nil
}

func (loan *paydayLoan) Amount() uint {
	return loan.amount
}

func (loan *paydayLoan) Term() Term {
	return loan.term
}
func (loan *paydayLoan) repay(amount uint) (err error) {
	if amount > loan.remaining {
		return ErrRepaymentAmountTooHigh
	}
	loan.remaining -= amount
	return nil
}

// Term is a number of days for which a loan was taken
type Term uint

// ErrClientAlreadyHasLoan is returned when Client already has active (unpaid) loan
var ErrClientAlreadyHasLoan = errors.New("client_already_has_loan")

// ErrAmountTooHigh is returned when Client applied for a loan with excessive amount
var ErrAmountTooHigh = AmountTooHighStruct{errors.New("amount_too_high"), maximumAmountForFirstLoan}

// AmountTooHighStruct is an error struct used by ErrAmountTooHigh with additional MaxAmount field indicating the maximum amount of a loan
type AmountTooHighStruct struct {
	error
	MaxAmount int
}

// ErrRepaymentAmountTooHigh is returned when Client tried to repay more than remaining amount of a loan
var ErrRepaymentAmountTooHigh = errors.New("repayment_amount_too_high")
