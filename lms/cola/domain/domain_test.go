package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	amount    uint = 10000000
	term           = Term(30)
	ktpNumber      = "3522582509010002"
)

func TestClientApplyForLoan(t *testing.T) {
	client := NewClient("", "", ktpNumber)
	err := client.ApplyForLoan(amount, term)
	t.Run("active loan should be assigned to client", func(t *testing.T) {
		assert.True(t, client.HasActiveLoan())
		loan := client.ActiveLoan()
		assert.Equal(t, loan.Amount(), amount)
		assert.Equal(t, loan.Remaining(), amount)
		assert.Equal(t, loan.Term(), term)
	})
	t.Run("error should be nil", func(t *testing.T) {
		assert.Nil(t, err)
	})
}

func TestClientApplyForLoanTwice(t *testing.T) {
	client := NewClient("", "", ktpNumber)
	client.ApplyForLoan(amount, term)
	err := client.ApplyForLoan(amount, term)
	t.Run("should return error", func(t *testing.T) {
		assert.Equal(t, "client_already_has_loan", err.Error())
		assert.Equal(t, err, ErrClientAlreadyHasLoan)
	})
}

func TestClientApplyForMoreThanMaxAmount(t *testing.T) {
	client := NewClient("", "", ktpNumber)
	err := client.ApplyForLoan(10000001, term)
	assert.Equal(t, err, ErrAmountTooHigh)
	assert.Equal(t, "amount_too_high", err.Error())
	assert.Equal(t, 1600, err.(AmountTooHighStruct).MaxAmount)
}

func TestClientRepaysLoanPart(t *testing.T) {
	client, loan := clientWithLoan(10000000)
	err := client.Repay(50)
	t.Run("Remaining amount should be 50", func(t *testing.T) {
		assert.Equal(t, uint(50), loan.Remaining())
	})
	t.Run("Should have active loan", func(t *testing.T) {
		assert.True(t, client.HasActiveLoan())
		assert.Equal(t, loan, client.ActiveLoan())
	})
	t.Run("Error should be nil", func(t *testing.T) {
		assert.Nil(t, err)
	})
}

func TestClientRepaysWholeLoan(t *testing.T) {
	client, loan := clientWithLoan(10000000)
	err := client.Repay(10000000)
	t.Run("Should not have active loan", func(t *testing.T) {
		assert.False(t, client.HasActiveLoan())
		assert.Nil(t, client.ActiveLoan())
	})
	t.Run("Remaning amount of repaid loan should 0", func(t *testing.T) {
		assert.Equal(t, loan.Remaining(), uint(0))
	})
	t.Run("error should be nil", func(t *testing.T) {
		assert.Nil(t, err)
	})
}

func TestClientRepaysTooMuch(t *testing.T) {
	client, _ := clientWithLoan(10000000)
	err := client.Repay(110)
	assert.Equal(t, ErrRepaymentAmountTooHigh, err)
	assert.Equal(t, "repayment_amount_too_high", err.Error())
}

func clientWithLoan(amount uint) (Client, Loan) {
	var client = NewClient("", "", ktpNumber)
	client.ApplyForLoan(amount, term)
	return client, client.ActiveLoan()
}
