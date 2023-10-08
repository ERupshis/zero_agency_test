// Package config server's setting parser. Applies flags and environments. Environments are prioritized.
package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env"
)

// Config server's settings.
type Config struct {
	DatabaseDSN      string
	DatabaseIdleConn int
	DatabaseOpenConn int
	Host             string
}

// Parse main func to parse variables.
func Parse() Config {
	var config = Config{}
	checkFlags(&config)
	checkEnvironments(&config)
	return config
}

// FLAGS PARSING.
const (
	flagAddress          = "a"
	flagDatabaseDSN      = "d"
	flagDatabaseIdleConn = "di"
	flagDatabaseOpenConn = "do"
)

// checkFlags checks flags of app's launch.
func checkFlags(config *Config) {
	// main app.
	flag.StringVar(&config.Host, flagAddress, "localhost:8080", "server endpoint")

	// postgres.
	flag.StringVar(&config.DatabaseDSN, flagDatabaseDSN, "postgres://postgres:postgres@localhost:5432/zero_agency_db?sslmode=disable", "database DSN")
	flag.IntVar(&config.DatabaseIdleConn, flagDatabaseIdleConn, 3, "database max idle connections")
	flag.IntVar(&config.DatabaseOpenConn, flagDatabaseOpenConn, 3, "database max open connections")

	flag.Parse()
}

// ENVIRONMENTS PARSING.
// envConfig struct of environments suitable for server.
type envConfig struct {
	DatabaseDSN      string `env:"DB_DSN"`
	DatabaseIdleConn string `env:"DB_MAX_IDLE_CONN"`
	DatabaseOpenConn string `env:"DB_MAX_OPEN_CONN"`
	Host             string `env:"ADDRESS"`
}

// checkEnvironments checks environments suitable for server.
func checkEnvironments(config *Config) {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		log.Fatal(err)
	}

	// main app.
	_ = SetEnvToParamIfNeed(&config.Host, envs.Host)

	// postgres.
	_ = SetEnvToParamIfNeed(&config.DatabaseDSN, envs.DatabaseDSN)
	_ = SetEnvToParamIfNeed(&config.DatabaseIdleConn, envs.DatabaseIdleConn)
	_ = SetEnvToParamIfNeed(&config.DatabaseOpenConn, envs.DatabaseOpenConn)
}
