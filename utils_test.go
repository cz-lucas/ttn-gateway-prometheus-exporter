package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyExistsInConfig(t *testing.T) {
	const testKey = "TEST_ENV_VAR"

	// Clean up environment after test
	defer os.Unsetenv(testKey)

	t.Run("Key does not exist", func(t *testing.T) {
		os.Unsetenv(testKey)
		assert.False(t, keyExistsInConfig(testKey))
	})

	t.Run("Key exists with non-empty value", func(t *testing.T) {
		os.Setenv(testKey, "somevalue")
		assert.True(t, keyExistsInConfig(testKey))
	})

	t.Run("Key exists but is empty", func(t *testing.T) {
		os.Setenv(testKey, "")
		assert.False(t, keyExistsInConfig(testKey))
	})
}

func TestGetEnvBool(t *testing.T) {
	const testKey = "BOOL_ENV_VAR"

	// Clean up environment after test
	defer os.Unsetenv(testKey)

	t.Run("Key does not exist", func(t *testing.T) {
		os.Unsetenv(testKey)
		result, err := getEnvBool(testKey, false)
		assert.False(t, result)
		assert.Nil(t, err)
	})

	t.Run("Key exists with true value", func(t *testing.T) {
		os.Setenv(testKey, "true")
		result, err := getEnvBool(testKey, false)
		assert.True(t, result)
		assert.Nil(t, err)
	})

	t.Run("Key exists with false value", func(t *testing.T) {
		os.Setenv(testKey, "false")
		result, err := getEnvBool(testKey, true)
		assert.False(t, result)
		assert.Nil(t, err)
	})

	t.Run("Key has a string", func(t *testing.T) {
		os.Setenv(testKey, "i am a string")
		result, err := getEnvBool(testKey, true)
		assert.False(t, result)
		expected := "invalid BOOL_ENV_VAR: strconv.ParseBool: parsing \"i am a string\": invalid syntax"
		assert.EqualError(t, err, expected)
	})
}

func TestGetEnvInt(t *testing.T) {
	const testKey = "INT_ENV_VAR"

	// Clean up environment after test
	defer os.Unsetenv(testKey)

	t.Run("Key does not exist", func(t *testing.T) {
		os.Unsetenv(testKey)
		result, err := getEnvInt(testKey, 1337)
		assert.Equal(t, result, 1337)
		assert.Nil(t, err)
	})

	t.Run("Key exists with value", func(t *testing.T) {
		os.Setenv(testKey, "42")
		result, err := getEnvInt(testKey, 1337)
		assert.Equal(t, result, 42)
		assert.Nil(t, err)
	})

	t.Run("Key has a string", func(t *testing.T) {
		os.Setenv(testKey, "i am a string")
		result, err := getEnvBool(testKey, true)
		if result != false {
			t.Errorf("Expected false, got true")
		}
		expected := "invalid INT_ENV_VAR: strconv.ParseBool: parsing \"i am a string\": invalid syntax"
		if err == nil || err.Error() != expected {
			t.Errorf("Unexpected error: got %v, want %v", err, expected)
		}
	})
}

func TestGetEnvString(t *testing.T) {
	const testKey = "STRING_ENV_VAR"

	// Clean up environment after test
	defer os.Unsetenv(testKey)

	t.Run("Key does not exist", func(t *testing.T) {
		os.Unsetenv(testKey)
		result := getEnvString(testKey, "hello")
		assert.Equal(t, result, "hello")
	})

	t.Run("Key exists with value", func(t *testing.T) {
		os.Setenv(testKey, "world")
		result := getEnvString(testKey, "world")
		assert.Equal(t, result, "world")
	})
}
