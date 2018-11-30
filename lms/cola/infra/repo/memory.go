// Package repo provides lms.ClientRepo implementation which stores everything in memory
package repo

import (
	"github.com/briyanadityatama/goLoans/lms/cola"
	"github.com/briyanadityatama/goLoans/lms/cola/domain"
)

type memoryClientRepo struct {
	clientsByKTPNumber map[string]domain.Client
}

// NewMemoryClientRepo returns a new instance of repository holding everything in memory
func NewMemoryClientRepo() cola.ClientRepo {
	clients := make(map[string]domain.Client)
	return &memoryClientRepo{clientsByKTPNumber: clients}
}

func (repo *memoryClientRepo) ByKTPNumber(ktplNumber string) (domain.Client, bool, error) {
	client, ok := repo.clientsByKTPNumber[ktplNumber]
	return client, ok, nil
}

func (repo *memoryClientRepo) Save(client domain.Client) error {
	ktplNumber := client.KTPNumber()
	repo.clientsByKTPNumber[ktplNumber] = client
	return nil
}
