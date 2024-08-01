package conf_test

import (
	"os"
	"strings"

	"github.com/fritzkeyzer/conf"
)

func ExampleMustLoad() {
	type Cfg struct {
		Debug bool   `flag:"--debug"`           // enable debug logs
		Host  string `flag:"--host" env:"HOST"` // host URL
		DB    struct {
			User string `env:"DB_USER" flag:"--db-user"`
			Pass string `env:"DB_PASS" flag:"--db-pass" secret:"db-pass"`
		}
	}

	// fake env vars and cli flags
	_ = os.Setenv("HOST", "localhost:1111")
	_ = os.Setenv("DB_USER", "not_the_root_user")
	fakeCliFlags := " --debug --host=localhost:8888" // notice how the flag here will override the env var
	os.Args = strings.Split(fakeCliFlags, " ")

	// one-liner loads the config
	cfg := conf.MustLoad[Cfg](conf.LoadCfg{Env: true, Flags: true})

	conf.Print(cfg)
	//Output:
	// -----------------------------------
	//   Debug     = true
	//   Host      = "localhost:8888"
	//   DB
	//     .User   = "not_the_root_user"
	//     .Pass   *** (len=0)
	// -----------------------------------
}
