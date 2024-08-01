package conf

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
)

const (
	maxPrintWidth = 90
	secretMask    = "***"
)

// Print wraps PrintToString and prints the result to stdout.
// Example output:
//
//	Host      = "localhost"
//	Verbose   = false
//	DB
//	  .Name   = "app"
//	  .User   = "user"
//	  .Pass   ***
func Print(ptr any) {
	_, _ = fmt.Fprintln(os.Stdout, PrintToString(ptr))
}

// PrintToString returns a string representation of the config struct. Secrets are masked.
// Example output:
//
//	Host      = "localhost"
//	Verbose   = false
//	DB
//	  .Name   = "app"
//	  .User   = "user"
//	  .Pass   ***
func PrintToString(ptr any) string {
	v := reflect.ValueOf(ptr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return "ERROR: config.Print: requires a struct as an argument"
	}

	fields := flattenFields(v, nil)

	buf := bytes.NewBuffer(nil)
	table := tablewriter.NewWriter(buf)
	table.SetHeaderLine(false)
	table.SetColWidth(maxPrintWidth)
	table.SetColumnSeparator(" ")
	table.SetCenterSeparator("-")

	for _, field := range fields {
		_, env := field.EnvVar()
		_, flag := field.FlagName()
		_, secret := field.SecretKey()

		printVal := true
		if field.value.Kind() == reflect.Struct {
			printVal = env || flag || secret
		}

		name := field.name
		if len(field.path) > 0 {
			name = strings.Repeat("  ", len(field.path)) + "." + name
		}

		value := ""
		if printVal {
			value = fmt.Sprintf("= %#v", field.value.Interface())
			if secret {
				value = secretMask

				valueLength := len(fmt.Sprint(field.value.Interface()))
				if field.value.Kind() == reflect.String {
					valueLength = len(field.value.String())
				} else if field.value.Kind() == reflect.Slice {
					valueLength = field.value.Len()
				}

				if valueLength == 0 {
					value = secretMask + " (len=0)"
				}
			}
		}

		if len(value) > maxPrintWidth-1 {
			value = value[:maxPrintWidth-1] + "..."
		}

		table.Append([]string{name, value})
	}

	table.Render()

	render := buf.String()
	lines := strings.Split(render, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}

	return strings.Join(lines, "\n")
}
