package main

import (
	"fmt"

	"github.com/fritzkeyzer/conf"
)

type Config struct {
	Host string `env:"HOST" flag:"--host"`
	DB   struct {
		Name string `env:"DB_NAME"`
		User string `env:"DB_USER" secret:"db-user"`
		Pass string `env:"DB_PASS" secret:"db-pass"`
	}
}

func main() {
	var cfg Config

	// create an implementation of the SecretSource interface
	// in reality this could be GCP SecretManager, AWS Secrets Manager, Vault, etc.
	secretManager := &SecretManager{}

	// LoadSecrets will load fields tagged as secret from the SecretManager
	if err := conf.LoadSecrets(&cfg, secretManager); err != nil {
		panic(err)
	}

	// the secret fields that are also tagged as env vars can be exported as env vars
	// for eg: to be saved to an .env file
	envFile := conf.ExportSecretsAsEnvVars(&cfg)
	fmt.Print(envFile) //nolint:forbidigo

	// Output:
	// export DB_USER="user from secret manager"
	// export DB_PASS="secret password 1337"
}

type SecretManager struct{}

func (sm *SecretManager) Load(key string) (string, bool, error) {
	if key == "db-user" {
		return "user from secret manager", true, nil
	}

	if key == "db-pass" {
		return "secret password 1337", true, nil
	}

	return "", false, nil
}
