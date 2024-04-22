package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var configTestCases = []struct {
	envVars     map[string]string
	expected    *Config
	expectedErr error
}{
	{
		envVars: map[string]string{
			"MAX_WORKERS":    "10",
			"MAX_QUEUE_SIZE": "100",
			"SERVER_ADDRESS": "localhost:8080",
			"READ_TIMEOUT":   "10s",
			"WRITE_TIMEOUT":  "10s",
		},
		expected: &Config{
			MaxWorkers:    10,
			MaxQueueSize:  100,
			ServerAddress: "localhost:8080",
			ReadTimeout:   10 * time.Second,
			WriteTimeout:  10 * time.Second,
		},
		expectedErr: nil,
	},
}

func TestGetConfig(t *testing.T) {
	for _, tc := range configTestCases {
		for key, value := range tc.envVars {
			os.Setenv(key, value)
			defer os.Unsetenv(key)
		}

		actual, err := GetConfig()

		assert.Equal(t, tc.expectedErr, err)
		assert.Equal(t, tc.expected, actual)
	}
}
