package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
}

func (c *MockDatabase) GetUserdataFor(username string) (userdata Userdata, err error) {

	if username == "joe" {
		userdata.HashedPassword = []byte("hashed_password")
		userdata.Roles = []byte("{admin}")
	}

	return
}

func TestGetUserdataFor_ExistingUser(t *testing.T) {
	// Set up the test environment
	mockDB := &MockDatabase{}

	// Define the expected user data
	expectedUserdata := Userdata{
		HashedPassword: []byte("hashed_password"),
		Roles:          []byte("{admin}"),
	}

	// Call the method to get user data
	userdata, err := mockDB.GetUserdataFor("joe")

	// Check for errors
	assert.NoError(t, err, "Expected no error")

	// Compare the returned user data with the expected values
	assert.Equal(t, expectedUserdata, userdata, "Returned user data does not match expected values")
}
