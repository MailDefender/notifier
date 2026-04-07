package utils

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func GetEnvString(varName string, defaultValue string) string {
	val := os.Getenv(varName)
	if val == "" {
		return defaultValue
	}

	return val
}

func GetEnvInt(varName string, defaultValue int) int {
	valStr := GetEnvString(varName, "")
	if valStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valStr)
	if err != nil {
		logrus.WithError(err).Warnf("invalid integer value for environment variable %s: %s, using default value %d", varName, valStr, defaultValue)
		return defaultValue
	}

	return value
}
