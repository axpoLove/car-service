package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config is a common config structure
type Config struct {
	Server         Server
	Postgres       Postgres
	CarInfoService CarInfoService
}

// NewConfig returns a new config instance
func NewConfig() (*Config, error) {
	// if we don't have .env file, it means that our app works in prod
	// so we don't check a error
	_ = godotenv.Load()
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load env variables: %w", err)
	}
	return &cfg, nil
}

// CarInfoService is car info service
type CarInfoService struct {
	URL            string        `envconfig:"CAR_INFO_SERVICE_URL" required:"true"`
	RequestTimeout time.Duration `envconfig:"CAR_INFO_REQUEST_TIMEOUT" default:"5s"`
}

// Postgres is a config for postgres database
type Postgres struct {
	Host          string `envconfig:"DB_HOST" required:"true"`
	Port          int    `envconfig:"DB_PORT" required:"true"`
	DB            string `envconfig:"POSTGRES_DB" required:"true"`
	User          string `envconfig:"POSTGRES_DB" required:"true"`
	Password      string `envconfig:"POSTGRES_PASSWORD" required:"true"`
	MigrationsDir string `envconfig:"POSTGRES_MIGRATIONS_DIR" default:"migrations"`
}

// ConnString returns connection string
func (p *Postgres) ConnString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.Host,
		p.Port,
		p.User,
		p.Password,
		p.DB,
	)
}

// Server is a config for http servers
type Server struct {
	CarPort           int           `envconfig:"CAR_SERVER_PORT" default:"8080"`
	CarInfoPort       int           `envconfig:"CAR_INFO_SERVER_PORT" default:"8081"`
	ReadHeaderTimeout time.Duration `envconfig:"SERVER_READ_HEADER_TIMEOUT" default:"15s"`
}
