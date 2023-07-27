<!-- Code generated by gomarkdoc. DO NOT EDIT -->

[![Go Report Card](https://goreportcard.com/badge/github.com/fritzkeyzer/conf)](https://goreportcard.com/report/github.com/fritzkeyzer/conf)
[![GitHub Actions](https://github.com/fritzkeyzer/conf/workflows/Test/badge.svg)](https://github.com/fritzkeyzer/conf/actions?query=workflow%3ATest+branch%3Amaster)
[![Go Reference](https://pkg.go.dev/badge/github.com/fritzkeyzer/conf.svg)](https://pkg.go.dev/github.com/fritzkeyzer/conf)

# conf

```go
import "github.com/fritzkeyzer/conf"
```

Package conf provides a set of utilities for mapping configuration settings \(from env vars, flags and secret managers\) to struct fields.

## Index

- [func ExportSecretsAsEnvVars\(ptr any\) string](<#ExportSecretsAsEnvVars>)
- [func GetFlag\(flag string, args \[\]string\) \(val string, found bool\)](<#GetFlag>)
- [func LoadEnv\(ptr any\) error](<#LoadEnv>)
- [func LoadFlags\(ptr any\) error](<#LoadFlags>)
- [func LoadSecrets\(ptr any, source SecretSource\) error](<#LoadSecrets>)
- [func Print\(ptr any\)](<#Print>)
- [func PrintToString\(ptr any\) string](<#PrintToString>)
- [type SecretSource](<#SecretSource>)


<a name="ExportSecretsAsEnvVars"></a>
## func [ExportSecretsAsEnvVars](<https://github.com/fritzkeyzer/conf/blob/master/secret.go#L61>)

```go
func ExportSecretsAsEnvVars(ptr any) string
```

ExportSecretsAsEnvVars recursively scans the struct fields to find fields that have both the secret and env tags. ExportSecretsAsEnvVars returns a string containing a list of export statements \- using the env var specified by the tag and the value of the field.

<details><summary>Example</summary>
<p>



```go
package main

import (
	"fmt"

	"github.com/fritzkeyzer/conf"
)

func main() {
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

}
```

#### Output

```
export DB_USER='user'
export DB_PASS='pass'
```

</p>
</details>

<a name="GetFlag"></a>
## func [GetFlag](<https://github.com/fritzkeyzer/conf/blob/master/flag.go#L25>)

```go
func GetFlag(flag string, args []string) (val string, found bool)
```

GetFlag is a utility to extract a flag from a slice of CLI args. It returns the value of the flag and a boolean indicating whether the flag was found. For example, args could be os.Args\[1:\]. flag should include the prefix, eg: "\-\-verbose" or "\-v" GetFlag supports the following formats:

```
flag=value
flag="value"
flag='value'
flag value
flag "value"
flag 'value'
```

<details><summary>Example</summary>
<p>



```go
package main

import (
	"fmt"

	"github.com/fritzkeyzer/conf"
)

func main() {
	args := []string{"nonsense", "--xyz=abc", "nonsense", "-v"}

	xyz, _ := conf.GetFlag("--xyz", args)
	_, verbose := conf.GetFlag("-v", args)

	fmt.Printf("xyz = %q, verbose = %v", xyz, verbose)

}
```

#### Output

```
xyz = "abc", verbose = true
```

</p>
</details>

<a name="LoadEnv"></a>
## func [LoadEnv](<https://github.com/fritzkeyzer/conf/blob/master/env.go#L16>)

```go
func LoadEnv(ptr any) error
```

LoadEnv recursively scans struct fields for the env tag then sets the values from the corresponding env var. Eg:

```
type Config struct {
	Host string `env:"HOST"`
}
```

<a name="LoadFlags"></a>
## func [LoadFlags](<https://github.com/fritzkeyzer/conf/blob/master/flag.go#L65>)

```go
func LoadFlags(ptr any) error
```

LoadFlags recursively scans struct fields for the flag tag then sets the values from CLI flags. Eg:

```
type Config struct {
	Host    string `flag:"--host"`
	Verbose bool   `flag:"-v"`
}
```

<a name="LoadSecrets"></a>
## func [LoadSecrets](<https://github.com/fritzkeyzer/conf/blob/master/secret.go#L23>)

```go
func LoadSecrets(ptr any, source SecretSource) error
```

LoadSecrets recursively scans struct fields for the secret tag then sets the values from the secret SecretSource. Eg:

```
type Config struct {
	Host string `secret:"host"`
}
```

<a name="Print"></a>
## func [Print](<https://github.com/fritzkeyzer/conf/blob/master/print.go#L18>)

```go
func Print(ptr any)
```

Print wraps PrintToString and prints the result to stdout.

<a name="PrintToString"></a>
## func [PrintToString](<https://github.com/fritzkeyzer/conf/blob/master/print.go#L23>)

```go
func PrintToString(ptr any) string
```

PrintToString returns a string representation of the config struct. Secrets are masked.

<a name="SecretSource"></a>
## type [SecretSource](<https://github.com/fritzkeyzer/conf/blob/master/secret.go#L11-L15>)

SecretSource interface allows any secret manager to be used, by wrapping it in a type that implements this interface.

```go
type SecretSource interface {
    // Load a secret from the source. Returns the secret value, a boolean indicating if the secret was found and an error.
    // NOTE: Load should not return an error if the secret was not found, but should instead return "", false, nil.
    Load(key string) (string, bool, error)
}
```

Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)