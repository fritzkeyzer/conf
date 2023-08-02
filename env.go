package conf

import (
	"fmt"
	"os"
)

// LoadEnv recursively scans struct fields for the env tag then sets the values from the corresponding env var.
// Eg:
//
//	type Config struct {
//		Host string `env:"HOST"`
//	}
func LoadEnv(ptr any) error {
	fields, err := FlattenStructFields(ptr)
	if err != nil {
		return err
	}

	for _, field := range fields {
		envVar, env := field.EnvVar()
		if !env {
			continue
		}

		envVal, found := os.LookupEnv(envVar)

		if err := field.SetValue(envVal, found); err != nil {
			return fmt.Errorf("failed to set field %q from env var: %w", field.field.Name, err)
		}
	}

	return nil
}
