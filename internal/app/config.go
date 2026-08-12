package app

import (
	"crypto/rsa"
	"log/slog"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type Config struct {
	dbUser string
	dbPass string
	dbHost string
	dbPort string
	dbName string

	privateRSAKey *rsa.PrivateKey
	publicRSAKey  *rsa.PublicKey
}

func (c *Config) DBUser() string { return c.dbUser }
func (c *Config) DBPass() string { return c.dbPass }
func (c *Config) DBHost() string { return c.dbHost }
func (c *Config) DBPort() string { return c.dbPort }
func (c *Config) DBName() string { return c.dbName }

func (c *Config) PrivateRSAKey() *rsa.PrivateKey { return c.privateRSAKey }
func (c *Config) PublicRSAKey() *rsa.PublicKey   { return c.publicRSAKey }

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Info(".env не найден, используются системные переменные окружения")
	}

	privateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	publicKeyPath := os.Getenv("PUBLIC_KEY_PATH")

	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	return &Config{
		dbUser: os.Getenv("DB_USER"),
		dbPass: os.Getenv("DB_PASS"),
		dbHost: os.Getenv("DB_HOST"),
		dbPort: os.Getenv("DB_PORT"),
		dbName: os.Getenv("DB_NAME"),

		privateRSAKey: privateKey,
		publicRSAKey:  publicKey,
	}, nil
}
