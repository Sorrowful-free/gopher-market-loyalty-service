package config

import (
	"flag"
	"os"
)

type LocalConfig struct {
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
	jwtSecret            string
	isProduction         bool
	isStackTrace         bool
}

func NewLocalConfig() Config {

	return &LocalConfig{
		runAddress:           "",
		databaseURI:          "",
		accrualSystemAddress: "",
		jwtSecret:            "",
		isProduction:         false,
		isStackTrace:         false,
	}
}

func (c *LocalConfig) Parse() error {

	flag.StringVar(&c.runAddress, "a", "localhost:8080", "run address")
	flag.StringVar(&c.databaseURI, "d", "localhost:5432", "database URI")
	flag.StringVar(&c.accrualSystemAddress, "r", "http://localhost:8080", "accrual system address")
	flag.StringVar(&c.jwtSecret, "j", "secret", "jwt secret")
	flag.BoolVar(&c.isProduction, "p", false, "is production")
	flag.BoolVar(&c.isStackTrace, "s", false, "is stack trace")
	flag.Parse()

	runAddress := os.Getenv("RUN_ADDRESS")
	if runAddress != "" {
		c.runAddress = runAddress
	}

	databaseURI := os.Getenv("DATABASE_URI")
	if databaseURI != "" {
		c.databaseURI = databaseURI
	}
	accrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if accrualSystemAddress != "" {
		c.accrualSystemAddress = accrualSystemAddress
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		c.jwtSecret = jwtSecret
	}

	isProduction := os.Getenv("IS_PRODUCTION")
	if isProduction != "" {
		c.isProduction = isProduction == "true"
	}

	isStackTrace := os.Getenv("IS_STACK_TRACE")
	if isStackTrace != "" {
		c.isStackTrace = isStackTrace == "true"
	}

	return nil
}

func (c *LocalConfig) RunAddress() string {
	return c.runAddress
}

func (c *LocalConfig) DatabaseURI() string {
	return c.databaseURI
}

func (c *LocalConfig) AccrualSystemAddress() string {
	return c.accrualSystemAddress
}

func (c *LocalConfig) JwtSecret() string {
	return c.jwtSecret
}

func (c *LocalConfig) IsProduction() bool {
	return c.isProduction
}

func (c *LocalConfig) IsStackTrace() bool {
	return c.isStackTrace
}
