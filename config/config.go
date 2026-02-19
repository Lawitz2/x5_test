package config

import (
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/url"
	"os"
	"strconv"
)

type DBConfig struct {
	Host       string
	Port       string
	User       string
	Password   string
	Name       string
	SslMode    string
	ConnString string
}

type Config struct {
	HTTPPort      string
	GRPCPort      string
	MigrationsDir string
	PageLimit     int
	DBConfig      DBConfig
}

// NewConfig загружает конфигурацию из .env файла или переменных окружения.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	dbCfg := DBConfig{}

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file: %s", err.Error())
	}

	cfg.HTTPPort = os.Getenv("HTTP_PORT")
	cfg.GRPCPort = os.Getenv("GRPC_PORT")
	cfg.MigrationsDir = os.Getenv("MIGRATIONS_DIR")

	pageLimit := os.Getenv("PAGE_LIMIT")
	cfg.PageLimit, err = strconv.Atoi(pageLimit)
	if err != nil {
		log.Printf("Error parsing PAGE_LIMIT environment variable: %s", err.Error())
		cfg.PageLimit = 100
	}

	dbCfg.Host = os.Getenv("DB_HOST")
	dbCfg.Port = os.Getenv("DB_PORT")
	dbCfg.User = os.Getenv("DB_USER")
	dbCfg.Password = os.Getenv("DB_PASSWORD")
	dbCfg.Name = os.Getenv("DB_NAME")
	dbCfg.SslMode = os.Getenv("DB_SSL_MODE")

	if dbCfg.Port == "" {
		dbCfg.Port = "5432"
	}
	if dbCfg.SslMode == "" {
		dbCfg.SslMode = "disable"
	}

	query := url.Values{}
	query.Add("sslmode", dbCfg.SslMode)
	dbUrl := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(dbCfg.User, dbCfg.Password),
		Host:     net.JoinHostPort(dbCfg.Host, dbCfg.Port),
		Path:     dbCfg.Name,
		RawQuery: query.Encode(),
	}

	dbCfg.ConnString = dbUrl.String()

	cfg.DBConfig = dbCfg

	return cfg, nil
}
