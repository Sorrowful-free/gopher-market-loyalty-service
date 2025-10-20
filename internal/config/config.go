package config

type Config interface {
	Parse() error

	RunAddress() string
	DatabaseURI() string
	AccrualSystemAddress() string
	JwtSecret() string
}
