package conf_test

import (
	"fmt"

	"github.com/fritzkeyzer/conf"
)

func ExampleExportSecretsAsEnvVars() {
	type Config struct {
		DBUser string `env:"DB_USER" secret:"db-user"`
		DBPass string `env:"DB_PASS" secret:"db-pass"`
	}

	// create a config with the secret fields filled in
	var cfg Config
	cfg.DBUser = "user"
	cfg.DBPass = "pass"

	// print the config as env vars
	// NOTE: only the fields that are tagged with a secret name and env var are exported
	fmt.Println(conf.ExportSecretsAsEnvVars(&cfg))

	// Output:
	// export DB_USER='user'
	// export DB_PASS='pass'
}
