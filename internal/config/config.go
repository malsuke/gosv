package config

/**
 * 構造体の書き方は以下を参照
 * @link https://github.com/caarlos0/env
 */
type Config struct {
	ENV_GITHUB_PAT string `env:"GITHUB_PAT"`
}
