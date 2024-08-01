package conf

type LoadCfg struct {
	Env           bool
	Flags         bool
	SecretsLoader SecretsLoader
}

// Load config from multiple sources.
// T should be a struct with tagged fields:
//   - secrets: `secret:mySecretValue`
//   - env vars: `env:MY_ENV_VAR`
//   - CLI flags: `flag:--flag`
//
// Sources are loaded in the following order:
//  1. First load secrets from SecretsLoader (if not nil)
//  2. Then environment variables - which will override secrets
//  3. Finally command line flags - which override both secrets and env vars
func Load[T any](cfg LoadCfg) (T, error) {
	var v T
	if cfg.SecretsLoader != nil {
		err := LoadSecrets(&v, cfg.SecretsLoader)
		if err != nil {
			return v, err
		}
	}
	if cfg.Env {
		err := LoadEnv(&v)
		if err != nil {
			return v, err
		}
	}
	if cfg.Flags {
		err := LoadFlags(&v)
		if err != nil {
			return v, err
		}
	}
	return v, nil
}

// MustLoad is a wrapper for Load which will panic if Load returns an error.
func MustLoad[T any](cfg LoadCfg) T {
	v, err := Load[T](cfg)
	if err != nil {
		panic(err)
	}
	return v
}
