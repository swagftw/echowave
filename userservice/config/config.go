package config

import (
	"log/slog"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type (
	// Config is the configuration for the gateway.
	Config struct {
		Port     int       `yaml:"port"`
		Debug    bool      `yaml:"debug"`
		OTel     *OTel     `yaml:"otel"`
		Postgres *Postgres `yaml:"postgres"`
	}

	OTel struct {
		DSN           string `yaml:"dsn"`
		Service       string `yaml:"service"`
		Version       string `yaml:"version"`
		DeploymentEnv string `yaml:"deploymentEnv"`
	}

	Postgres struct {
		URL string `yaml:"url"`
	}
)

// ReadConfig reads the config from the given path.
func ReadConfig(configPath string) (*Config, error) {
	byteArray, err := os.ReadFile(configPath)
	if err != nil {
		slog.Error("error reading config file", "path", configPath, "error", err)

		return nil, err
	}

	config := new(Config)

	err = yaml.Unmarshal(byteArray, config)
	if err != nil {
		slog.Error("error unmarshalling config", "error", err)

		return nil, err
	}

	err = config.readFromEnv()

	return config, err
}

// readFromEnv reads the config from environment variables.
func (config *Config) readFromEnv() error {
	if port := os.Getenv("PORT"); port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			slog.Error("error parsing port", "error", err)

			return err
		}

		config.Port = p
	}

	if debug := os.Getenv("DEBUG"); debug != "" {
		d, err := strconv.ParseBool(debug)
		if err != nil {
			slog.Error("error parsing debug", "error", err)

			return err
		}

		config.Debug = d
	}

	if dsn := os.Getenv("UPTRACE_DSN"); dsn != "" {
		config.OTel.DSN = dsn
	}

	if postgresURL := os.Getenv("POSTGRES_URL"); postgresURL != "" {
		config.Postgres.URL = postgresURL
	}

	return nil
}
