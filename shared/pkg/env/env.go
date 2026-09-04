package env

import "github.com/caarlos0/env/v11"

// LoadEnvConfiguration parses environment variables into the given struct T.
// Fields must be tagged with `env` and optionally `envDefault`.
func LoadEnvConfiguration[T any]() (*T, error) {
	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
