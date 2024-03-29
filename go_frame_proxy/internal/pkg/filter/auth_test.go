package filter

import (
	"testing"

	"github.com/leobowenwang/go_frame_proxy/internal/pkg/database"
	"github.com/stretchr/testify/assert"
)

type MockDatabase struct{}

func (c *MockDatabase) GetUserdataFor(username string) (userdata database.Userdata, err error) {

	if username == "joe" {
		userdata.HashedPassword = []byte("$2a$12$9TKMv5j6rigr.o/c98JoP.ffNvMgAVNVqaVeNA3tG4SQdLhjyBUju")
		userdata.Roles = []byte("{admin}")
	}

	return
}

func TestAuthentication_CorrectCredentials(t *testing.T) {
	Username := "joe"
	Password := "joemama"

	mockDatabase := &MockDatabase{}

	isAuthenticated, _, _ := BCryptAuthenticate(mockDatabase, Username, Password)

	assert.True(t, isAuthenticated)
}

func TestAuthentication_WrongUsername(t *testing.T) {
	Username := "joey"
	Password := "joemama"

	mockDatabase := &MockDatabase{}

	isAuthenticated, _, _ := BCryptAuthenticate(mockDatabase, Username, Password)

	assert.False(t, isAuthenticated)

}

func TestAuthentication_WrongPassword(t *testing.T) {
	Username := "joe"
	Password := "joemami"

	mockDatabase := &MockDatabase{}

	isAuthenticated, _, _ := BCryptAuthenticate(mockDatabase, Username, Password)

	assert.False(t, isAuthenticated)

}
