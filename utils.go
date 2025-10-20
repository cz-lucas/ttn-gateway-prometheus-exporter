package main

import (
	"fmt"
	"os"
	"strconv"
)

func keyExistsInConfig(name string) bool {
	val := os.Getenv(name)
	if val == "" {
		return false
	}
	return true
}

func getEnvBool(name string, defaultVal bool) (bool, error) {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal, nil
	}
	i, err := strconv.ParseBool(val)
	if err != nil {
		return false, fmt.Errorf("invalid %s: %w", name, err)
	}
	return i, nil
}

func getEnvInt(name string, defaultVal int) (int, error) {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal, nil
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", name, err)
	}
	return i, nil
}

func getEnvString(name string, defaultVal string) string {
	val := os.Getenv(name)
	if val == "" {
		return defaultVal
	}
	return val
}
