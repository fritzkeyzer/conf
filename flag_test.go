package conf

import (
	"fmt"
	"testing"
)

func ExampleGetFlag() {
	args := []string{"nonsense", "--xyz=abc", "nonsense", "-v"}

	xyz, _ := GetFlag("--xyz", args)
	_, verbose := GetFlag("-v", args)

	fmt.Printf("xyz = %q, verbose = %v", xyz, verbose)

	// Output: xyz = "abc", verbose = true
}

func TestGetFlag(t *testing.T) {
	type args struct {
		flag string
		args []string
	}
	tests := []struct {
		name      string
		args      args
		wantVal   string
		wantFound bool
	}{
		// TODO: Add test cases.
		{
			name: "no args",
			args: args{
				flag: "--xyz",
				args: []string{},
			},
			wantVal:   "",
			wantFound: false,
		},
		{
			name: "=",
			args: args{
				flag: "--xyz",
				args: []string{"nonsense", "--xyz=abc", "nonsense"},
			},
			wantVal:   "abc",
			wantFound: true,
		},
		{
			name: "=quoted '",
			args: args{
				flag: "--xyz",
				args: []string{"nonsense", "--xyz='abc'", "nonsense"},
			},
			wantVal:   "abc",
			wantFound: true,
		},
		{
			name: "=quoted \"",
			args: args{
				flag: "--xyz",
				args: []string{"nonsense", `--xyz="abc"`, "nonsense"},
			},
			wantVal:   "abc",
			wantFound: true,
		},
		{
			name: "space",
			args: args{
				flag: "--xyz",
				args: []string{"nonsense", "--xyz", "abc", "nonsense"},
			},
			wantVal:   "abc",
			wantFound: true,
		},
		{
			name: "space quoted \"",
			args: args{
				flag: "--xyz",
				args: []string{"nonsense", "--xyz", `"abc"`, "nonsense"},
			},
			wantVal:   "abc",
			wantFound: true,
		},
		{
			name: "space quoted '",
			args: args{
				flag: "--xyz",
				args: []string{"nonsense", "--xyz", `'abc'`, "nonsense"},
			},
			wantVal:   "abc",
			wantFound: true,
		},
		{
			name: "overlap",
			args: args{
				flag: "--xy",
				args: []string{"nonsense", "--xyz=def", "nonsense"},
			},
			wantVal:   "",
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVal, gotFound := GetFlag(tt.args.flag, tt.args.args)
			if gotVal != tt.wantVal {
				t.Errorf("GetFlag() gotVal = %v, want %v", gotVal, tt.wantVal)
			}
			if gotFound != tt.wantFound {
				t.Errorf("GetFlag() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}
