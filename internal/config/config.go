package config

import (
	"runtime"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"

	"github.com/tarkovskynik/Golang-ninja-project/internal/domain"
	"github.com/tarkovskynik/Golang-ninja-project/pkg/logger"
)

func parseConfigFile(configDir string) error {
	viper.AddConfigPath(configDir)
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return nil
}

func unmarshal(cfg *domain.Config) error {
	if err := viper.UnmarshalKey("serverListener.tcp", &cfg.Server); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}
	if err := viper.UnmarshalKey("file", &cfg.File); err != nil {
		return err
	}
	return nil
}

func setFromEnv(cfg *domain.Config) error {

	if runtime.GOOS == "windows" {
		// windows specific code here...
		err := godotenv.Load()
		if err != nil {
			logger.Fatalf("Error loading .env file")
		}
	}

	if err := envconfig.Process("db", &cfg.Postgres); err != nil {
		return err
	}

	if err := envconfig.Process("auth", &cfg.Auth); err != nil {
		return err
	}
	return nil
}

func Init(configDir string) (*domain.Config, error) {
	viper.SetConfigName("config")
	if err := parseConfigFile(configDir); err != nil {
		return nil, err
	}

	var cfg domain.Config
	if err := unmarshal(&cfg); err != nil {
		return nil, err
	}
	if err := setFromEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func init() {
	// err := godotenv.Load()
	// if err != nil {
	// 	logrus.Fatal("Error loading .env file")
	// }
}
