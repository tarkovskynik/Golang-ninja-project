package domain

import (
	"time"

	"github.com/tarkovskynik/Golang-ninja-project/pkg/s3"
)

type NetServerConfig struct {
	Host string
	Port int
}

type PostgresConfig struct {
	Host     string `envconfig:"Host"`
	Port     int    `envconfig:"Port"`
	Username string `envconfig:"Username"`
	Name     string `envconfig:"Name"`
	SSLMode  string `envconfig:"SSLMode"`
	Password string `envconfig:"Password"`
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	Salt            string
	Secret          string
}

type FileConfig struct {
	Storage       s3.ConfigFileStorage
	MaxUploadSize int64                  // 10 megabytes = 10 << 20
	CheckTypes    map[string]interface{} // "image/jpeg": nil, "image/png": nil, ...
	Types         []string
}

type Config struct {
	Server   NetServerConfig
	Postgres PostgresConfig
	Auth     AuthConfig
	File     FileConfig
}
