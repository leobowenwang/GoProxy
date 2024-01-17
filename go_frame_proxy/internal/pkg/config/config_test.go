package config

import (
	"os"
	"path"
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func getWorkingDirectory() (pathURL string) {
	userMachine, _ := os.UserHomeDir()
	fmt.Printf("userMachine-variable: %s", userMachine)
	pathURL = path.Join(userMachine, "/work/sesame-gateway/sesame-gateway/internal/pkg/config/")
	return
}

func TestLoadConfig(t *testing.T) {

	// Test case 1: Test loading the valid config

	filePath := path.Join(getWorkingDirectory(), "testdata/test_valid_config.yml")

	config, err := LoadConfig(filePath)
	assert.NoError(t, err, "Expected no error for valid config")
	assert.Equal(t, 3000, config.Server.Port, "Unexpected value for server port")
	assert.Equal(t, "./cert/server.crt", config.Server.CertFile, "Unexpected value for certFile")
	assert.Equal(t, "./cert/server.key", config.Server.KeyFile, "Unexpected value for keyFile")
	assert.Equal(t, "sesame.postgres.database.azure.com", config.Database.Host, "Unexpected value for pgHost")
	assert.Equal(t, "userDB", config.Database.Database, "Unexpected value for pgDatabase")
	assert.Equal(t, "user", config.Database.User, "Unexpected value for pgUser")
	assert.Equal(t, "Tiniwo18", config.Database.Password, "Unexpected value for pgPassword")
	assert.Len(t, config.Proxy, 3, "Unexpected number of proxy configurations")
	assert.Equal(t, "/get", config.Proxy[0].Path, "Unexpected value for proxy[0].path")
	assert.Equal(t, "localhost:1000", config.Proxy[0].Host, "Unexpected value for proxy[0].host")
	assert.Equal(t, []string{"GET"}, config.Proxy[0].Methods, "Unexpected value for proxy[0].methods")
	assert.Equal(t, "/post", config.Proxy[1].Path, "Unexpected value for proxy[1].path")
	assert.Equal(t, "localhost:2000", config.Proxy[1].Host, "Unexpected value for proxy[1].host")
	assert.Equal(t, []string{"POST"}, config.Proxy[1].Methods, "Unexpected value for proxy[1].methods")
	assert.Equal(t, "/yes", config.Proxy[2].Path, "Unexpected value for proxy[2].path")
	assert.Equal(t, "localhost:2000", config.Proxy[2].Host, "Unexpected value for proxy[2].host")

	// Test case 2: Test loading the invalid config file
	filePath = path.Join(getWorkingDirectory(), "./testdata/test_invalid_config.yml")
	_, err = LoadConfig(filePath)
	assert.NoError(t, err, "Expected no error for valid config")
	assert.Equal(t, 3000, config.Server.Port, "Unexpected value for server port")
	assert.Equal(t, "./cert/server.crt", config.Server.CertFile, "Unexpected value for certFile")
	assert.Equal(t, "./cert/server.key", config.Server.KeyFile, "Unexpected value for keyFile")
	assert.Equal(t, "sesame.postgres.database.azure.com", config.Database.Host, "Unexpected value for pgHost")
	assert.Equal(t, "userDB", config.Database.Database, "Unexpected value for pgDatabase")
	assert.Equal(t, "user", config.Database.User, "Unexpected value for pgUser")
	assert.Equal(t, "Tiniwo18", config.Database.Password, "Unexpected value for pgPassword")
	assert.Len(t, config.Proxy, 3, "Unexpected number of proxy configurations")
	assert.Equal(t, "/get", config.Proxy[0].Path, "Unexpected value for proxy[0].path")
	assert.Equal(t, "localhost:1000", config.Proxy[0].Host, "Unexpected value for proxy[0].host")
	assert.Equal(t, []string{"GET"}, config.Proxy[0].Methods, "Unexpected value for proxy[0].methods")
	assert.Equal(t, "/post", config.Proxy[1].Path, "Unexpected value for proxy[1].path")
	assert.Equal(t, "localhost:2000", config.Proxy[1].Host, "Unexpected value for proxy[1].host")
	assert.Equal(t, []string{"POST"}, config.Proxy[1].Methods, "Unexpected value for proxy[1].methods")
	assert.Equal(t, "/yes", config.Proxy[2].Path, "Unexpected value for proxy[2].path")
	assert.Equal(t, "localhost:2000", config.Proxy[2].Host, "Unexpected value for proxy[2].host")

}

func TestReadConfigFile(t *testing.T) {

	// Test case 1: Test reading the valid config file
	filePath := path.Join(getWorkingDirectory(), "/testdata/test_valid_config.yml")
	config, err := readConfigFile(filePath)
	assert.NoError(t, err, "Expected no error for valid config")
	assert.Equal(t, 3000, config.Server.Port, "Unexpected value for server port")
	assert.Equal(t, "./cert/server.crt", config.Server.CertFile, "Unexpected value for certFile")
	assert.Equal(t, "./cert/server.key", config.Server.KeyFile, "Unexpected value for keyFile")
	assert.Equal(t, "sesame.postgres.database.azure.com", config.Database.Host, "Unexpected value for pgHost")
	assert.Equal(t, "userDB", config.Database.Database, "Unexpected value for pgDatabase")
	assert.Equal(t, "user", config.Database.User, "Unexpected value for pgUser")
	assert.Equal(t, "Tiniwo18", config.Database.Password, "Unexpected value for pgPassword")
	assert.Len(t, config.Proxy, 3, "Unexpected number of proxy configurations")
	assert.Equal(t, "/get", config.Proxy[0].Path, "Unexpected value for proxy[0].path")
	assert.Equal(t, "localhost:1000", config.Proxy[0].Host, "Unexpected value for proxy[0].host")
	assert.Equal(t, []string{"GET"}, config.Proxy[0].Methods, "Unexpected value for proxy[0].methods")
	assert.Equal(t, "/post", config.Proxy[1].Path, "Unexpected value for proxy[1].path")
	assert.Equal(t, "localhost:2000", config.Proxy[1].Host, "Unexpected value for proxy[1].host")
	assert.Equal(t, []string{"POST"}, config.Proxy[1].Methods, "Unexpected value for proxy[1].methods")
	assert.Equal(t, "/yes", config.Proxy[2].Path, "Unexpected value for proxy[2].path")
	assert.Equal(t, "localhost:2000", config.Proxy[2].Host, "Unexpected value for proxy[2].host")

	// Test reading the invalid config file
	filePath = path.Join(getWorkingDirectory(), "/testdata/test_invalid_read_config.txt")
	_, err = readConfigFile(filePath)
	assert.Error(t, err, "Expected error for invalid config file")

}
