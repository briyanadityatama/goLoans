package cola

import "github.com/briyanadityatama/goLoans/lms/cola/domain"

type fakeClientRepo struct {
	clientsByKTPNumber map[string]domain.Client
}

// NewFakeClientRepo returns ClientRepo fake implementation storing everything in memory which is useful for testing lms.Lms without real database
func NewFakeClientRepo() ClientRepo {
	return &fakeClientRepo{clientsByKTPNumber: make(map[string]domain.Client)}
}

func (repo *fakeClientRepo) ByKTPNumber(ktpNumber string) (domain.Client, bool, error) {
	client, ok := repo.clientsByKTPNumber[ktpNumber]
	return client, ok, nil
}

func (repo *fakeClientRepo) Save(client domain.Client) error {
	ktpNumber := client.KTPNumber()
	repo.clientsByKTPNumber[ktpNumber] = client
	return nil
}
