package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zapcore"
)

type mockConfig struct {
	mock.Mock
}

func (m *mockConfig) GetString(key string) string {
	args := m.Called(key)

	return args.String(0)
}

func (m *mockConfig) GetInt(key string) int {
	args := m.Called(key)

	return args.Int(0)
}

func TestSetup(t *testing.T) {
	tests := []struct {
		logLevel string
		expected zapcore.Level
	}{
		{"debug", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"error", zapcore.ErrorLevel},
		{"invalid", zapcore.InfoLevel},
	}

	for _, tt := range tests {
		t.Run(tt.logLevel, func(t *testing.T) {
			mockCfg := new(mockConfig)
			mockCfg.On("GetString", "log_level").Return(tt.logLevel)

			Setup(mockCfg)
			assert.NotNil(t, l, "logger should be initialized")

			l = nil
		})
	}
}
