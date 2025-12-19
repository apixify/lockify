package cmd

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

type mockClearUseCase struct {
	executeFunc func(ctx context.Context) error
	executed    bool
}

func (m *mockClearUseCase) Execute(ctx context.Context) error {
	m.executed = true
	if m.executeFunc != nil {
		return m.executeFunc(ctx)
	}
	return nil
}

func TestClearCommand_Success(t *testing.T) {
	mockUseCase := &mockClearUseCase{}
	buildUseCase := func() app.ClearCachedPassphraseUc {
		return mockUseCase
	}
	mockLogger := &test.MockLogger{}

	cmd := NewClearCommand(buildUseCase, mockLogger)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.True(t, mockUseCase.executed)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.SuccessLogs)
}

func TestClearCommand_UseCaseError(t *testing.T) {
	errMsg := "execute failed"
	mockUseCase := &mockClearUseCase{
		executeFunc: func(ctx context.Context) error {
			return fmt.Errorf("%s", errMsg)
		},
	}
	buildUseCase := func() app.ClearCachedPassphraseUc {
		return mockUseCase
	}
	mockLogger := &test.MockLogger{}

	cmd := NewClearCommand(buildUseCase, mockLogger)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, errMsg, err.Error())
	assert.True(t, mockUseCase.executed)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
	assert.Count(t, 1, mockLogger.ErrorLogs)
}
