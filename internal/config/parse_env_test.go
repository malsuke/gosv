package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseConfigFromEnv(t *testing.T) {
	t.Run("success from env", func(t *testing.T) {
		t.Setenv("GITHUB_PAT", "sample_pat_from_env")

		cfg := &Config{}
		err := cfg.ParseConfigFromEnv()

		assert.NoError(t, err)
		assert.Equal(t, "sample_pat_from_env", cfg.ENV_GITHUB_PAT)
	})

	t.Run("no env var set", func(t *testing.T) {
		t.Setenv("GITHUB_PAT", "")

		cfg := &Config{}
		err := cfg.ParseConfigFromEnv()

		assert.NoError(t, err)
		assert.Empty(t, cfg.ENV_GITHUB_PAT)
	})
}
