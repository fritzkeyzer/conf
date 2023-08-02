package conf

import (
	"fmt"
)

// SecretSource interface allows any secret manager to be used, by wrapping it in a type that implements this interface.
type SecretSource interface {
	// Load a secret from the source. Returns the secret value, a boolean indicating if the secret was found and an error.
	// NOTE: Load should not return an error if the secret was not found, but should instead return "", false, nil.
	Load(key string) ([]byte, bool, error)
}

// LoadSecrets recursively scans struct fields for the secret tag then sets the values from the secret SecretSource.
// Eg:
//
//	type Config struct {
//		Host string `secret:"host"`
//	}
func LoadSecrets(ptr any, source SecretSource) error {
	fields, err := FlattenStructFields(ptr)
	if err != nil {
		return err
	}

	for _, field := range fields {
		secretKey, secret := field.SecretKey()
		if !secret {
			continue
		}

		val, found, err := source.Load(secretKey)
		if err != nil {
			return fmt.Errorf("failed to load secret %q: %w", secretKey, err)
		}

		if err := field.setBytes(val, found); err != nil {
			return fmt.Errorf("failed to set field %q from secret source: %w", field.field.Name, err)
		}
	}

	return nil
}
