package main

import (
	"os"

	"github.com/fritzkeyzer/conf"
)

type Config struct {
	Host    string `env:"HOST" flag:"--host"`
	Verbose bool   `flag:"-v"`
	Count   int    `flag:"--count"`
	DB      struct {
		Name string `env:"DB_NAME"`
		User string `env:"DB_USER" secret:"db-user"`
		Pass string `env:"DB_PASS" secret:"db-pass"`
	}
}

// main demonstrates various functions of the conf package
//   - LoadEnv loads fields from environment variables
//   - LoadFlags loads fields from command line flags
//   - LoadSecrets loads fields from a secret manager
//   - Print prints the config to stdout
func main() {
	// for demo purposes, we set the env vars here
	os.Setenv("HOST", "localhost")
	os.Setenv("DB_NAME", "app")
	os.Setenv("DB_USER", "user from env")
	os.Setenv("DB_PASS", "pass from env")

	var cfg Config

	if err := conf.LoadEnv(&cfg); err != nil {
		panic(err)
	}

	// fields can be overridden by flags, eg: host, verbose or count
	if err := conf.LoadFlags(&cfg); err != nil {
		panic(err)
	}

	// fields can be loaded from a secret manager
	if err := conf.LoadSecrets(&cfg, &SecretManager{}); err != nil {
		panic(err)
	}

	// notice how the secret fields are masked with ***
	conf.Print(&cfg)

	// Output:
	// ---------------------------
	//  Host      = "localhost"
	//  Verbose   = false
	//  DB
	//    .Name   = "app"
	//    .User   = "user"
	//    .Pass   ***
	// ---------------------------
}

// SecretManager is a mock secret manager, for demo purposes
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
