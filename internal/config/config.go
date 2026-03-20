package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	APIToken  string `mapstructure:"api-token"`
	AccountID string `mapstructure:"account-id"`
}

func Load() *Config {
	viper.SetConfigName(".mailtrap")
	viper.SetConfigType("yaml")

	home, err := os.UserHomeDir()
	if err == nil {
		viper.AddConfigPath(home)
	}
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("MAILTRAP")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()

	var cfg Config
	_ = viper.Unmarshal(&cfg)
	return &cfg
}

func RequireAccountID() (string, error) {
	id := viper.GetString("account-id")
	if id == "" {
		return "", fmt.Errorf("account-id is required: set via --account-id, MAILTRAP_ACCOUNT_ID, or ~/.mailtrap.yaml")
	}
	return id, nil
}

func RequireAPIToken() (string, error) {
	token := viper.GetString("api-token")
	if token == "" {
		return "", fmt.Errorf("api-token is required: set via --api-token, MAILTRAP_API_TOKEN, or ~/.mailtrap.yaml")
	}
	return token, nil
}

func ConfigFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".mailtrap.yaml"
	}
	return filepath.Join(home, ".mailtrap.yaml")
}
