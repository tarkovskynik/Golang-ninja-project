package domain

import "time"

type NetServerConfig struct {
	Host string
	Port int
}

type PostgresConfig struct {
	Host     string
	Port     int
	Username string
	Name     string
	SSLMode  string
	Password string
}

type AuthConfig struct {
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl"`
	Salt            string
	Secret          string
}

type FileConfig struct {
	Size       int
	Extensions []string
}

type Config struct {
	Server   NetServerConfig
	Postgres PostgresConfig
	Auth     AuthConfig
	File     FileConfig
}
