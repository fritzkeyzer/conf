package conf

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

const flagPrefix = "-"

// GetFlag is a utility to extract a flag from a slice of CLI args.
// It returns the value of the flag and a boolean indicating whether the flag was found.
// For example, args could be os.Args[1:].
// flag should include the prefix, eg: "--verbose" or "-v"
// GetFlag supports the following formats:
//
//	flag=value
//	flag="value"
//	flag='value'
//	flag value
//	flag "value"
//	flag 'value'
func GetFlag(flag string, args []string) (val string, found bool) {
	for i := range args {
		if !strings.HasPrefix(args[i], flag) {
			continue
		}

		if strings.HasPrefix(args[i], flag+"=") {
			val = strings.TrimPrefix(args[i], flag+"=")
			val = strings.Trim(val, `"`)
			val = strings.Trim(val, `'`)
			return val, true
		}

		if args[i] == flag {
			// if there are more args and the next arg is not a flag
			if len(args) > i+1 && !strings.HasPrefix(args[i+1], flagPrefix) {
				// the next arg is the value
				val = args[i+1]
				val = strings.Trim(val, `"`)
				val = strings.Trim(val, `'`)

				return val, true
			}

			// else return found, without any value (for example boolean flags)
			return "", true
		}
	}

	// all args have been checked and the flag has not been found
	return "", false
}

// LoadFlags recursively scans struct fields for the flag tag then sets the values from CLI flags.
// Eg:
//
//	type Config struct {
//		Host    string `flag:"--host"`
//		Verbose bool   `flag:"-v"`
//	}
func LoadFlags(ptr any) error {
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
		flagName, flag := field.flagName()
		if !flag {
			continue
		}

		flagVar, found := GetFlag(flagName, os.Args[1:])
		if !found {
			continue
		}

		if err := field.setValue(flagVar, found); err != nil {
			return fmt.Errorf("failed to set field %q from flag: %w", field.field.Name, err)
		}
	}

	return nil
}
