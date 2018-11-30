package lms

type fakeLms struct {
	clientsByKTPNumber map[string]Client
}

// NewFakeLms returns Lms fake implementation storing everything in memory which is useful for testing GUI and other clients (such as REST server)
func NewFakeLms() Lms {
	return &fakeLms{clientsByKTPNumber: make(map[string]Client)}
}

func (lms *fakeLms) RegisterClient(clientData ClientData) (Client, error) {
	newClient := fakeClient{clientData.Gender, clientData.KTPNumber, clientData.BirthDate, clientData.Name}
	lms.clientsByKTPNumber[clientData.KTPNumber] = newClient
	return newClient, nil
}

func (lms *fakeLms) ClientByKTPNumber(ktpNumber string) (client Client, found bool, err error) {
	client, ok := lms.clientsByKTPNumber[ktpNumber]
	return client, ok, nil
}

func (lms *fakeLms) ApplyForLoan(ktpNumber string, amount uint, term uint) error {
	panic("implement me")
}

type fakeClient struct {
	gender, ktpNumber, birthDate, name string
}

func (client fakeClient) Gender() string {
	return client.gender
}

func (client fakeClient) KTPNumber() string {
	return client.ktpNumber
}

func (client fakeClient) BirthDate() string {
	return client.birthDate
}

func (client fakeClient) Name() string {
	return client.name
}

func (client fakeClient) HasActiveLoan() bool {
	panic("implement me")
}
