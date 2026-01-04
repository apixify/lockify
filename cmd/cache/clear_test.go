package cache

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

type mockClearUseCase struct {
	executeFunc func(vctx *model.VaultContext) error
	executed    bool
}

func (m *mockClearUseCase) Execute(vctx *model.VaultContext) error {
	m.executed = true
	if m.executeFunc != nil {
		return m.executeFunc(vctx)
	}
	return nil
}

func TestClearCommand_Success(t *testing.T) {
	mockClearUseCase := &mockClearUseCase{}
	mockPassphrase := &test.MockPassphraseService{}
	mockClearEnvUseCase := app.NewClearEnvCachedPassphraseUseCase(mockPassphrase)
	mockLogger := &test.MockLogger{}

	cmd := NewClearCommand(
		mockClearUseCase,
		mockClearEnvUseCase,
		mockLogger,
		cli.NewCommandContext(),
	)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.True(t, mockClearUseCase.executed)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.SuccessLogs)
}

func TestClearCommand_UseCaseError(t *testing.T) {
	mockClearUseCase := &mockClearUseCase{
		executeFunc: func(vctx *model.VaultContext) error {
			return fmt.Errorf("%s", test.ErrMsgExecuteFailed)
		},
	}
	mockPassphrase := &test.MockPassphraseService{}

	mockClearEnvUseCase := app.NewClearEnvCachedPassphraseUseCase(mockPassphrase)
	mockLogger := &test.MockLogger{}

	cmd := NewClearCommand(
		mockClearUseCase,
		mockClearEnvUseCase,
		mockLogger,
		cli.NewCommandContext(),
	)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, test.ErrMsgExecuteFailed, err.Error())
	assert.True(t, mockClearUseCase.executed)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
	assert.Count(t, 1, mockLogger.ErrorLogs)
}
