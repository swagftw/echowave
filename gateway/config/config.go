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
		ReadTimeout  int          `yaml:"readTimeout"`
		WriteTimeout int          `yaml:"writeTimeout"`
		IdleTimeout  int          `yaml:"idleTimeout"`
		Port         int          `yaml:"port"`
		Debug        bool         `yaml:"debug"`
		OTel         *OTel        `yaml:"otel"`
		UserService  *UserService `yaml:"userService"`
	}

	OTel struct {
		DSN           string `yaml:"dsn"`
		Service       string `yaml:"service"`
		Version       string `yaml:"version"`
		DeploymentEnv string `yaml:"deploymentEnv"`
	}

	UserService struct {
		GrpcURL string `yaml:"grpcURL"`
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
	if readTimeout := os.Getenv("READ_TIMEOUT"); readTimeout != "" {
		rt, err := strconv.Atoi(readTimeout)
		if err != nil {
			slog.Error("error parsing read timeout", "error", err)

			return err
		}

		config.ReadTimeout = rt
	}

	if writeTimeout := os.Getenv("WRITE_TIMEOUT"); writeTimeout != "" {
		wt, err := strconv.Atoi(writeTimeout)
		if err != nil {
			slog.Error("error parsing write timeout", "error", err)

			return err
		}

		config.WriteTimeout = wt
	}

	if idleTimeout := os.Getenv("IDLE_TIMEOUT"); idleTimeout != "" {
		it, err := strconv.Atoi(idleTimeout)
		if err != nil {
			slog.Error("error parsing idle timeout", "error", err)

			return err
		}

		config.IdleTimeout = it
	}

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

	return nil
}
