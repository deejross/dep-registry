package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/deejross/dep-registry/util"
)

const envPrefix = "GOREG_"

// Config object.
type Config struct {
	BinStorePath  string
	MetaStorePath string
	SigningKey    string
	Port          string
}

// FromFile gets a Config object from a file.
func FromFile(c *Config, name string) (*Config, error) {
	if c == nil {
		c = &Config{}
	}

	f, err := os.OpenFile(name, os.O_RDONLY, 0600)
	if err != nil {
		return c, err
	}

	if err := json.NewDecoder(f).Decode(c); err != nil {
		return c, err
	}

	return c, nil
}

// FromEnvironment gets a Config object from environment variables.
func FromEnvironment(c *Config) *Config {
	if c == nil {
		c = &Config{}
	}

	if v := os.Getenv(envPrefix + "BINSTORE_PATH"); len(v) > 0 {
		c.BinStorePath = v
	}
	if v := os.Getenv(envPrefix + "METASTORE_PATH"); len(v) > 0 {
		c.MetaStorePath = v
	}
	if v := os.Getenv(envPrefix + "SIGNING_KEY"); len(v) > 0 {
		c.SigningKey = v
	}
	if v := os.Getenv(envPrefix + "PORT"); len(v) > 0 {
		c.Port = v
	}

	return c
}

// Validate configuration.
func (c *Config) Validate() error {
	if len(c.BinStorePath) == 0 {
		c.BinStorePath = "boltdb://binstore.bolt"
	}
	if len(c.MetaStorePath) == 0 {
		c.MetaStorePath = "boltdb://metastore.bolt"
	}
	if len(c.SigningKey) == 0 {
		log.Println("WARNING: No signing key specified, generating a temporary key")
		c.SigningKey = util.UUID4()
	}
	if len(c.Port) == 0 {
		c.Port = "8080"
	}

	return nil
}
