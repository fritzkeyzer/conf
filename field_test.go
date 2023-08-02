package conf

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_flattenFields(t *testing.T) {
	type Config struct {
		Host    string `env:"HOST"`
		private string // private fields are skipped
		Port    int    `env:"PORT"`
		DB      struct {
			User string `env:"DB_USER"`
		}
		Data   []byte `env:"DATA"`
		Nested struct {
			Deep string // skipped since parent is tagged with env
		} `env:"NESTED"`
	}

	cfg := Config{}

	v := reflect.ValueOf(&cfg)
	v = v.Elem()

	want := []Field{
		{
			path:  nil,
			name:  "Host",
			field: v.Type().Field(0),
			value: v.Field(0),
		},
		{
			path:  nil,
			name:  "Port",
			field: v.Type().Field(2),
			value: v.Field(2),
		},
		{
			path:  nil,
			name:  "DB",
			field: v.Type().Field(3),
			value: v.Field(3),
		},
		{
			path:  []string{"DB"},
			name:  "User",
			field: v.Field(3).Type().Field(0),
			value: v.Field(3).Field(0),
		},
		{
			path:  nil,
			name:  "Data",
			field: v.Type().Field(4),
			value: v.Field(4),
		},
		{
			path:  nil,
			name:  "Nested",
			field: v.Type().Field(5),
			value: v.Field(5),
		},
	}

	got := flattenFields(v, nil)

	if len(got) != len(want) {
		t.Fatalf("len(got) != len(want): %v != %v", len(got), len(want))
	}

	for i := range got {
		fieldsEqual(t, got[i], want[i])
	}
}

func fieldsEqual(t *testing.T, got, want Field) {
	gotPath := strings.Join(got.path, ".")
	wantPath := strings.Join(want.path, ".")
	if gotPath != wantPath {
		t.Errorf("got.path != want.path: %v != %v", gotPath, wantPath)
	}

	if got.name != want.name {
		t.Errorf("got.name != want.name: %v != %v", got.name, want.name)
	}

	if got.field.Name != want.field.Name {
		t.Errorf("got.field.Name != want.field.Name: %v != %v", got.field.Name, want.field.Name)
	}

	if got.field.Tag != want.field.Tag {
		t.Errorf("got.field.Tag != want.field.Tag: %v != %v", got.field.Tag, want.field.Tag)
	}

	if got.value.String() != want.value.String() {
		t.Errorf("got.value != want.value: %+v != %+v", got.value.Interface(), want.value.Interface())
	}
}

func Test_envEndToEnd(t *testing.T) {
	type Config struct {
		String string `env:"TEST_STRING"`
		Int    int    `env:"TEST_INT"`
		Bool   bool   `env:"TEST_BOOL"`
		Struct struct {
			String string `env:"TEST_STRUCT_STRING"`
		}
		Bytes []byte `env:"TEST_BYTES"`
	}

	// initialise a config struct with some values
	want := Config{}
	want.String = "string"
	want.Int = 42
	want.Bool = true
	want.Struct.String = `string ' with spaces " quotes
		and newlines`
	want.Bytes = []byte("placeholder bytes from a string")

	// flatten fields
	fields, err := FlattenStructFields(&want)
	if err != nil {
		t.Fatalf("FlattenStructFields: %v", err)
	}

	// export env vars
	for _, f := range fields {
		envVar, env := f.EnvVar()
		if !env {
			continue
		}

		envVal, err := f.ExportValue()
		if err != nil {
			t.Errorf("ExportValue: %v", err)
		}

		err = os.Setenv(envVar, envVal)
		if err != nil {
			t.Errorf("os.Setenv: %v", err)
		}

		t.Logf("export %s=%s", envVar, envVal)
	}

	// create 2nd config and load from env vars
	var got Config
	if err := LoadEnv(&got); err != nil {
		t.Fatalf("LoadEnv: %v", err)
	}

	wantString := fmt.Sprintf("%+v", want)
	gotString := fmt.Sprintf("%+v", got)

	if wantString != gotString {
		t.Fatalf("got != want:\ngot:\n%v\nwant:\n%v", gotString, wantString)
	}
}
