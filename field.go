package conf

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

const (
	envTag    = "env"
	flagTag   = "flag"
	secretTag = "secret"
)

// Field represents a struct field
type Field struct {
	path []string
	name string

	field reflect.StructField
	value reflect.Value
}

// FlattenStructFields returns a flat slice of Field from recursively traversing the struct fields of v.
//   - unexported fields are omitted
//   - fields marked with an env, flag or secret tag are included, but their children are not
func FlattenStructFields(ptr any) ([]Field, error) {
	v := reflect.ValueOf(ptr)
	if v.Kind() != reflect.Ptr {
		return nil, errors.New("requires a pointer to struct")
	}

	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return nil, errors.New("requires a pointer to struct")
	}

	return flattenFields(v, nil), nil
}

func flattenFields(v reflect.Value, path []string) []Field {
	t := v.Type()

	var fields []Field
	for i := 0; i < t.NumField(); i++ {
		// skip unexported fields
		if !t.Field(i).IsExported() {
			continue
		}

		f := Field{
			path:  path,
			name:  t.Field(i).Name,
			field: t.Field(i),
			value: v.Field(i),
		}

		fields = append(fields, f)

		// do not recurse into fields that have the env, flag or secret tags
		_, env := f.EnvVar()
		_, flag := f.FlagName()
		_, secret := f.SecretKey()
		if env || flag || secret {
			continue
		}

		if f.field.Type.Kind() == reflect.Struct {
			subFields := flattenFields(f.value, append(path, f.name))
			fields = append(fields, subFields...)
		}
	}

	return fields
}

// EnvVar returns the `env` tag value and a bool indicating if the field has the `env` tag.
func (f *Field) EnvVar() (string, bool) {
	envVar := f.field.Tag.Get(envTag)
	if envVar != "" {
		return envVar, true
	}

	return "", false
}

// FlagName returns the `flag` tag value and a bool indicating if the field has the `flag` tag.
func (f *Field) FlagName() (string, bool) {
	flagName := f.field.Tag.Get(flagTag)
	if flagName != "" {
		return flagName, true
	}

	return "", false
}

// SecretKey returns the `secret` tag value and a bool indicating if the field has the `secret` tag.
func (f *Field) SecretKey() (string, bool) {
	secretKey := f.field.Tag.Get(secretTag)
	if secretKey != "" {
		return secretKey, true
	}

	return "", false
}

// ExportValue returns the value of the field as a string.
//   - []byte fields are base64 encoded
//   - string fields are not pre-processed
//   - all other types marshalled to JSON
func (f *Field) ExportValue() (string, error) {
	if f.value.Kind() == reflect.Slice && f.value.Type().Elem().Kind() == reflect.Uint8 {
		return base64.StdEncoding.EncodeToString(f.value.Bytes()), nil
	}

	if f.value.Kind() == reflect.String {
		return f.value.String(), nil
	}

	buf, err := json.Marshal(f.value.Interface())
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

// setString sets the underlying field value from a string.
//   - []byte fields are assumed to be base64 encoded
//   - string fields are not pre-processed
//   - all other types are assumed to be JSON encoded
func (f *Field) setString(rawVal string, found bool) error {
	if f.value.Kind() == reflect.Slice && f.value.Type().Elem().Kind() == reflect.Uint8 {
		if !found {
			return nil
		}

		dst := make([]byte, base64.StdEncoding.DecodedLen(len(rawVal)))
		n, err := base64.StdEncoding.Decode(dst, []byte(rawVal))
		if err != nil {
			return fmt.Errorf("decoding base64: %w", err)
		}

		if n > 0 {
			f.value.SetBytes(dst[:n])
		}

		return nil
	}

	switch f.value.Kind() {
	case reflect.Bool:
		if found && rawVal == "" {
			f.value.SetBool(true)
		} else if found && rawVal != "" {
			f.value.SetBool(rawVal == "true")
		}

	case reflect.String:
		if !found {
			return nil
		}
		f.value.Set(reflect.ValueOf(rawVal))

	default:
		if !found {
			return nil
		}

		val := f.value.Addr().Interface()
		if err := json.Unmarshal([]byte(rawVal), val); err != nil {
			return fmt.Errorf("%w, raw value: %q", err, rawVal)
		}
	}

	return nil
}
