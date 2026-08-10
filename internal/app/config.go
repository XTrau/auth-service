package app

import "os"

type Config struct {
	dbUser string
	dbPass string
	dbHost string
	dbPort string
	dbName string
}

func (c *Config) DBUser() string { return c.dbUser }
func (c *Config) DBPass() string { return c.dbPass }
func (c *Config) DBHost() string { return c.dbHost }
func (c *Config) DBPort() string { return c.dbPort }
func (c *Config) DBName() string { return c.dbName }

func LoadConfig() *Config {
	return &Config{
		dbUser: os.Getenv("DB_USER"),
		dbPass: os.Getenv("DB_PASS"),
		dbHost: os.Getenv("DB_HOST"),
		dbPort: os.Getenv("DB_PORT"),
		dbName: os.Getenv("DB_NAME"),
	}
}
