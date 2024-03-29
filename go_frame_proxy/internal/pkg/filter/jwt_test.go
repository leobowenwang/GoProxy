package filter

import (
	"encoding/base64"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/leobowenwang/go_frame_proxy/internal/pkg/util"
	"github.com/stretchr/testify/assert"
)

func getWorkingDirectory() (pathURL string) {
	userMachine, _ := os.UserHomeDir()
	pathURL = path.Join(userMachine, "/work/sesame-gateway/sesame-gateway/internal/pkg/filter")
	return
}

func TestCertfileCheck_WhenGivenCorrectFiles(t *testing.T) {

	cert, key, err := CheckCertFile(path.Join(getWorkingDirectory(), "/Test Certfiles/server.crt"), path.Join(getWorkingDirectory(), "/Test Certfiles/server.key"))

	if err != nil {
		assert.Fail(t, err.Error())
	}

	if cert == nil {
		assert.Fail(t, "certificate could not be created.")
	}

	if key == nil {
		assert.Fail(t, "privatekey could not be created.")
	}

}

func TestFileCheck_WhenGivenWrongFilePaths(t *testing.T) {

	certFile := path.Join(getWorkingDirectory(), "/Test Certfiles/server.crt")
	keyFile := path.Join(getWorkingDirectory(), "/Test Certfiles/server.key")

	_, _, err := CheckCertFile("Wrong Path", keyFile)
	util.Logger.Debug(err.Error())
	assert.True(t, strings.Contains(err.Error(), "error while reading the content of certificate."))

	_, _, err = CheckCertFile(certFile, "Wrong Path")
	util.Logger.Debug(err.Error())
	assert.True(t, strings.Contains(err.Error(), "error while reading the content of Keyfile."))

}

func TestGenerateJWT(t *testing.T) {

	username := "TestUsername"
	roles := "{admin}"

	privateKey, err := GetPrivateKey(path.Join(getWorkingDirectory(), "/Test Certfiles/server.key"))

	if err != nil {
		assert.Fail(t, "Error while trying to create private key")
	}

	token, err := GenerateJWT(username, roles, privateKey)

	if err != nil {
		assert.Fail(t, "Error while creating JWT")
	}

	tokenParts := strings.Split(token, ".")

	header := `{"alg":"RS256","typ":"JWT"}`

	bytes := []byte(header)
	encoded := base64.StdEncoding.EncodeToString(bytes)

	if tokenParts[0] != encoded {
		assert.Fail(t, "The JWT Header was not generated correctly.")
	}

}
