package config

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/leobowenwang/go_frame_proxy/internal/pkg/util"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Port     int    `yaml:"port" validate:"nonzero"`
		CertFile string `yaml:"certFile"`
		KeyFile  string `yaml:"keyFile"`
	} `yaml:"server"`
	Database struct {
		Host     string `yaml:"pgHost" validate:"nonzero"`
		Database string `yaml:"pgDatabase" validate:"nonzero"`
		User     string `yaml:"pgUser" validate:"nonzero"`
		Password string `yaml:"pgPassword" validate:"nonzero"`
		Timeout  int    `yaml:"pgTimeout" validate:"nonzero"`
	} `yaml:"postgres"`
	Proxy []struct {
		Path    string   `yaml:"path"`
		Host    string   `yaml:"host"`
		Methods []string `yaml:"methods"`
	} `yaml:"proxy"`
}

// LoadConfig looks for config files at the provided locations, if no locations are provided it uses the default location
func LoadConfig(location string) (Config, error) {

	filePaths := []string{"./config.yml", "./config.yaml"}
	if location != "" {
		filePaths = []string{location}
	}

	var config Config
	var err error
	config, err = readConfigFile(location)
	for _, file := range filePaths {

		if err == nil {
			config, err = readConfigFile(file)
			return config, nil
		}
		util.Logger.Error("Failed to read config file", zap.String("filePath", file), zap.Error(err))
	}

	return Config{}, fmt.Errorf("failed to read config files")
}

func readConfigFile(filePath string) (Config, error) {

	file, err := os.ReadFile(filePath)
	if err != nil {
		util.Logger.Error("Failed to read config file", zap.String("filePath", filePath), zap.Error(err))
		return Config{}, fmt.Errorf("failed to read file: %w", err)
	}

	var config Config
	decoder := yaml.NewDecoder(bytes.NewReader(file))
	for {
		err := decoder.Decode(&config)
		if err != nil {
			if err == io.EOF {
				break
			}
			util.Logger.Error("Failed to decode config file", zap.String("filePath", filePath), zap.Error(err))
			return Config{}, fmt.Errorf("failed to decode config file: %w", err)
		}
	}

	return config, nil
}
