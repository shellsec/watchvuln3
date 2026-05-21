package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvBool(t *testing.T) {
	assert := require.New(t)
	key := "WATCHVULN_TEST_ENV_BOOL"
	defer os.Unsetenv(key)

	assert.False(assertEnvBool(key))

	os.Setenv(key, "true")
	v, ok := EnvBool(key)
	assert.True(ok)
	assert.True(v)

	os.Setenv(key, "false")
	v, ok = EnvBool(key)
	assert.True(ok)
	assert.False(v)

	os.Setenv(key, "1")
	v, ok = EnvBool(key)
	assert.True(ok)
	assert.True(v)

	os.Unsetenv(key)
	_, ok = EnvBool(key)
	assert.False(ok)
}

func assertEnvBool(key string) bool {
	_, ok := EnvBool(key)
	return !ok
}
