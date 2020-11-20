package storage

// MockClient fulfills the Client interface so that it can be used in place of a real storage Client
type MockClient struct {
	DoGet func(key string) (string, error)
}

// NewMock returns a MockClient for use in unit tests
func NewMock() *MockClient {

	return &MockClient{
		DoGet: func(key string) (string, error) {
			return "value", nil
		},
	}
}

func (m *MockClient) Get(key string) (string, error) {
	return m.DoGet(key)
}
