package conf_test

import (
	"testing"

	"github.com/fritzkeyzer/conf"
)

func TestPrintToString(t *testing.T) {
	type Config struct {
		privateField string
		Host         string
		DB           struct {
			Host string `secret:"db-host"`
			User string
			Pass string `secret:"db-pass"`
		}
		ServiceCreds struct {
			User string
			Pass string
		} `secret:"service-creds"`
	}

	var cfg Config
	cfg.privateField = "hello private" // private fields are skipped
	cfg.Host = "localhost"
	cfg.DB.Host = "" // intentionally left blank (see len=0 in output)
	cfg.DB.User = "admin"
	cfg.DB.Pass = "admin ;)"
	cfg.ServiceCreds.User = "admin" // ServiceCreds will be printed on one line
	cfg.ServiceCreds.Pass = "admin"

	got := conf.PrintToString(&cfg)
	want := `
--------------------------------
  Host           = "localhost"  
  DB                            
    .Host        *** (len=0)    
    .User        = "admin"      
    .Pass        ***            
  ServiceCreds   ***            
--------------------------------
`

	if got != want {
		t.Fatalf("got != want: got:\n%v\nwant:\n%v", got, want)
	}
}
