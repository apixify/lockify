package cmd

import (
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

func TestVersionCommand_Success(t *testing.T) {
	mockLogger := &test.MockLogger{}

	cmd := NewVersionCommand(mockLogger)

	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.Count(t, 1, mockLogger.OutputLogs)
	assert.Contains(t, "Lockify CLI", mockLogger.OutputLogs[0])
	assert.Contains(t, "v0.0.0", mockLogger.OutputLogs[0])
}
