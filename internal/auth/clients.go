package auth

type mockClientService struct{}

func NewMockClientService() *mockClientService {
	return &mockClientService{}
}

func (*mockClientService) ValidateUser(user, password string) bool {
	return user == password
}
