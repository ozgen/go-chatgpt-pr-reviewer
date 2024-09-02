package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	const envKey = "TEST_ENV_VAR"
	const defaultValue = "default"

	os.Setenv(envKey, "exists")
	assert.Equal(t, "exists", GetEnv(envKey, defaultValue))

	os.Unsetenv(envKey)
	assert.Equal(t, defaultValue, GetEnv(envKey, defaultValue))
}
