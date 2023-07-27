package conf

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type structField struct {
	path []string
	name string

	field reflect.StructField
	value reflect.Value
}

// flattenFields returns a flat slice of structField from recursively traversing the struct fields of v.
func flattenFields(v reflect.Value, path []string) []structField {
	t := v.Type()

	var fields []structField
	for i := 0; i < t.NumField(); i++ {
		var f structField
		f.field = t.Field(i)
		f.value = v.Field(i)
		f.path = path
		f.name = f.field.Name

		fields = append(fields, f)

		// do not recurse into fields that have the env, flag or secret tags
		_, env := f.envVar()
		_, flag := f.flagName()
		_, secret := f.secretKey()
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

func (f *structField) envVar() (string, bool) {
	envVar := f.field.Tag.Get("env")
	if envVar != "" {
		return envVar, true
	}

	return "", false
}

func (f *structField) flagName() (string, bool) {
	flagName := f.field.Tag.Get("flag")
	if flagName != "" {
		return flagName, true
	}

	return "", false
}

func (f *structField) secretKey() (string, bool) {
	secretKey := f.field.Tag.Get("secret")
	if secretKey != "" {
		return secretKey, true
	}

	return "", false
}

func (f *structField) setValue(rawVal string, found bool) error {
	switch f.value.Kind() { //nolint:exhaustive
	default:
		return fmt.Errorf("unsupported kind: %s", f.value.Kind())

	case reflect.String:
		f.value.Set(reflect.ValueOf(rawVal))

	case reflect.Bool:
		if rawVal != "" {
			f.value.SetBool(isTrue(rawVal))
		} else {
			f.value.Set(reflect.ValueOf(found))
		}

	case reflect.Int:
		intVal, err := strconv.Atoi(rawVal)
		if err != nil {
			return fmt.Errorf("parsing int: %w", err)
		}
		f.value.SetInt(int64(intVal))

	case reflect.Struct:
		val := f.value.Addr().Interface()
		if err := json.Unmarshal([]byte(rawVal), val); err != nil {
			return fmt.Errorf("parsing json: %w, %q", err, rawVal)
		}
	}

	return nil
}

func isTrue(s string) bool {
	switch strings.ToLower(s) {
	case "true", "y", "yes", "1":
		return true
	}

	return false
}
