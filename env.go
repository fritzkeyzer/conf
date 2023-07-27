package conf

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

// LoadEnv recursively scans struct fields for the env tag then sets the values from the corresponding env var.
// Eg:
//
//	type Config struct {
//		Host string `env:"HOST"`
//	}
func LoadEnv(ptr any) error {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		return errors.New("requires a pointer argument")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("requires a pointer to struct argument")
	}

	fields := flattenFields(v, nil)

	for _, field := range fields {
		envVar, env := field.envVar()
		if !env {
			continue
		}

		envVal, ok := os.LookupEnv(envVar)
		if !ok {
			continue
		}

		if err := field.setValue(envVal, true); err != nil {
			return fmt.Errorf("failed to set field %q from env var: %w", field.field.Name, err)
		}
	}

	return nil
}
