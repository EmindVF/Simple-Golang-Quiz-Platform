package config

import (
	"encoding/json"
	"fmt"
	"os"

	"quiz_platform/internal/misc/apperrors"
)

type (
	Config struct {
		App       App       `json:"app"`
		Database  Database  `json:"database"`
		TokenInfo TokenInfo `json:"token_info"`
	}

	TokenInfo struct {
		PublicKeyPath  string `json:"public_key_path"`
		PrivateKeyPath string `json:"private_key_path"`
		PublicKey      []byte `json:"-"`
		PrivateKey     []byte `json:"-"`
		ExpiresIn      int32  `json:"expires_in"`
		MaxAge         int32  `json:"max_age"`
	}

	App struct {
		Port int `json:"port"`
	}

	Database struct {
		Host           string `json:"host"`
		Port           int    `json:"port"`
		User           string `json:"user"`
		Password       string `json:"password"`
		DBName         string `json:"dbname"`
		SSLMode        string `json:"sslmode"`
		TimeZone       string `json:"timezone"`
		InitScriptPath string `json:"init_script_path"`
	}
)

var GlobalConfig *Config

func ReadGlobalConfig(path string) error {
	GlobalConfig = &Config{}

	if path == "" {
		path = "./config/config.json"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return &apperrors.ErrInternal{
			Message: fmt.Sprintf("cannot read config file: %v", err.Error()),
		}
	}

	err = json.Unmarshal(data, GlobalConfig)
	if err != nil {
		return &apperrors.ErrInternal{
			Message: fmt.Sprintf("cannot parse config file: %v", err.Error())}
	}

	GlobalConfig.TokenInfo.PublicKey, err = os.ReadFile(GlobalConfig.TokenInfo.PublicKeyPath)
	if err != nil {
		panic(fmt.Errorf("fatal error reading config files: %v", err))
	}
	GlobalConfig.TokenInfo.PrivateKey, err = os.ReadFile(GlobalConfig.TokenInfo.PrivateKeyPath)
	if err != nil {
		panic(fmt.Errorf("fatal error reading config files: %v", err))
	}

	return nil
}
