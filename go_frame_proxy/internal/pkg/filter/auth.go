package filter

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/project-sesame/sesame-gateway/internal/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func BCryptAuthenticate(db database.IDatabase, username, password string) (bool, string, error) {

	var hashedPassword, roles []byte

	userdata, _ := db.GetUserdataFor(username)

	hashedPassword = userdata.HashedPassword
	roles = userdata.Roles

	// Check if the hash is a bcrypt hash and its length is at least 60 bytes
	if len(hashedPassword) < 60 || string(hashedPassword[0:2]) != "$2" {
		return false, "", echo.NewHTTPError(http.StatusForbidden, "HashedSecret is not a valid bcrypt hash, please reset password")
	}

	// Compare the plain-text password with the hashed password using bcrypt
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return false, "", nil
	}

	return true, string(roles), nil
}
