package main

import (
	"os"

	"github.com/fritzkeyzer/conf"
)

type Config struct {
	Host    string `env:"HOST" flag:"--host"`
	Verbose bool   `flag:"-v"`
	DB      struct {
		Name string `env:"DB_NAME"`
		User string `env:"DB_USER"`
		Pass string `env:"DB_PASS" secret:"db-pass"`
	}
}

func main() {
	// for demo purposes, we set the env vars here
	os.Setenv("HOST", "localhost")
	os.Setenv("DB_NAME", "app")
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASS", "a secret password")

	var cfg Config

	if err := conf.LoadEnv(&cfg); err != nil {
		panic(err)
	}

	// fields can be overridden by flags
	if err := conf.LoadFlags(&cfg); err != nil {
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
