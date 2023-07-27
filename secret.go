package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// SecretSource interface allows any secret manager to be used, by wrapping it in a type that implements this interface.
type SecretSource interface {
	// Load a secret from the source. Returns the secret value, a boolean indicating if the secret was found and an error.
	// NOTE: Load should not return an error if the secret was not found, but should instead return "", false, nil.
	Load(key string) (string, bool, error)
}

// LoadSecrets recursively scans struct fields for the secret tag then sets the values from the secret SecretSource.
// Eg:
//
//	type Config struct {
//		Host string `secret:"host"`
//	}
func LoadSecrets(ptr any, source SecretSource) error {
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
		secretKey, secret := field.secretKey()
		if !secret {
			continue
		}

		val, found, err := source.Load(secretKey)
		if err != nil {
			return fmt.Errorf("failed to load secret %q: %w", secretKey, err)
		}
		if !found {
			continue
		}

		if err := field.setValue(val, true); err != nil {
			return fmt.Errorf("failed to set field %q from secret source: %w", field.field.Name, err)
		}
	}

	return nil
}

// ExportSecretsAsEnvVars recursively scans the struct fields to find fields that have both the secret and env tags.
// ExportSecretsAsEnvVars returns a string containing a list of export statements - using the env var specified by the tag
// and the value of the field.
func ExportSecretsAsEnvVars(ptr any) string {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		panic("requires a pointer argument")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		panic("requires a pointer to struct argument")
	}

	fields := flattenFields(v, nil)

	str := ""
	for _, field := range fields {
		_, secret := field.secretKey()
		if !secret {
			continue
		}

		envVar, env := field.envVar()
		if !env {
			continue
		}

		valStr := fmt.Sprint(field.value.Interface())

		// if the field is a struct, marshal it to json instead
		if field.value.Kind() == reflect.Struct {
			buf, err := json.Marshal(field.value.Interface())
			if err != nil {
				panic(err)
			}
			valStr = string(buf)
		}

		str += fmt.Sprintf("export %s='%s'\n", envVar, valStr)
	}

	return str
}
