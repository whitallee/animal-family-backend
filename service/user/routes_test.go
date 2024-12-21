package user

import (
	"testing"

	"github.com/whitallee/animal-family-backend/types"
)

func TestUserServiceHandlers(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)

	t.Run("should fail if the user payload is invalid") //START HERE
}

type mockUserStore struct {
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, nil
}
func (m *mockUserStore) GetUserById(id int) (*types.User, error) {
	return nil, nil
}
func (m *mockUserStore) CreateUser(types.User) error {
	return nil
}
