[![Go Report Card](https://goreportcard.com/badge/github.com/fritzkeyzer/conf)](https://goreportcard.com/report/github.com/fritzkeyzer/conf)
[![GitHub Actions](https://github.com/fritzkeyzer/conf/workflows/Test/badge.svg)](https://github.com/fritzkeyzer/conf/actions?query=workflow%3ATest+branch%3Amaster)
[![Go Reference](https://pkg.go.dev/badge/github.com/fritzkeyzer/conf.svg)](https://pkg.go.dev/github.com/fritzkeyzer/conf)

# conf

Package conf provides a set of utilities for mapping configuration settings from env vars, flags and secret managers to struct fields.

```go
import "github.com/fritzkeyzer/conf"
```

## Features
Load struct fields from environment variables
```go
var Config struct {
    DBConn string `env:"DB_CONN"`
}

_ = conf.LoadEnv(&Config)
```

Load struct fields from CLI flags
```go
var Config struct {
    DBConn string `flag:"db-conn"`
}

_ = conf.LoadFlags(&Config)
```

Load struct fields from a secret manager
```go
var Config struct {
    DBConn string `secret:"db-conn"`
}

_ = conf.LoadSecrets(&Config, secretMgr) // where secretMgr implements the SecretSource interface
```

## Utilities
Parse flags from []string, eg: os.Args
```go
args := []string{"nonsense", "--xyz=abc", "nonsense", "-v"}

// get value, found from []string
value, found := GetFlag("--xyz", args) // value = "abc", found = true

// eg:
_, verbose := GetFlag("-v", args) // verbose = true
```

Print a config to stdout
```go
type Config struct {
    Host    string `env:"HOST" flag:"--host"`
    Verbose bool   `flag:"-v"`
    Count   int    `flag:"--count"`
    DB      struct {
    Name string `env:"DB_NAME"`
    User string `env:"DB_USER" secret:"db-user"`
    Pass string `env:"DB_PASS" secret:"db-pass"`
}

conf.Print(&cfg) // notice how the secret fields are masked with ***

// Output:
// ---------------------------
//  Host      = "localhost"
//  Verbose   = false
//  DB
//    .Name   = "app"
//    .User   = "user"
//    .Pass   ***
// ---------------------------
```

Flatten struct fields and iterate over them, eg: to export to a .env file
```go
fields, _ := conf.FlattenStructFields(&cfg)

envFile := ""
for _, field := range fields {
    envVar, e := field.EnvVar()
    if !e {
        continue
    }
    
    envVal, _ := field.ExportValue() // exports values in a way that is suitable for env vars and for loading back into the struct
    envFile += "export " + envVar + "=" + envVal + "\n"
}
```

## Example
Example from [examples/main.go](examples/main.go)
```go
package main

import (
	"fmt"
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

	// other features include flattening the struct fields
	fields, _ := conf.FlattenStructFields(&cfg)

	// this way you can iterate over the fields and do something with them
	// eg: exporting the values to a .env file
	envFile := ""
	for _, field := range fields {
		envVar, e := field.EnvVar()
		if !e {
			continue
		}

		envVal, _ := field.ExportValue()
		envFile += "export " + envVar + "=" + envVal + "\n"
	}

	fmt.Println(envFile)
	// Output:
	//export HOST=localhost
	//export DB_NAME=app
	//export DB_USER=user from secret manager
	//export DB_PASS=secret password 1337
}

// SecretManager is a mock secret manager, for demo purposes
type SecretManager struct{}

func (sm *SecretManager) Load(key string) ([]byte, bool, error) {
	if key == "db-user" {
		return []byte("user from secret manager"), true, nil
	}

	if key == "db-pass" {
		return []byte("secret password 1337"), true, nil
	}

	return nil, false, nil
}

```
