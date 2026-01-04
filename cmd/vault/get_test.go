package vault

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

type mockGetUseCase struct {
	executeFunc func(vctx *model.VaultContext, key string) (string, error)
	recievedKey string
}

func (m *mockGetUseCase) Execute(vctx *model.VaultContext, key string) (string, error) {
	m.recievedKey = key
	if m.executeFunc != nil {
		return m.executeFunc(vctx, key)
	}

	return "test_value", nil
}

func TestGetCommand_Success(t *testing.T) {
	mockUseCase := &mockGetUseCase{}
	mockLogger := &test.MockLogger{}
	cmd, _ := NewGetCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	if err := cmd.Flags().Set("key", "test_key"); err != nil {
		t.Fatalf("failed to set key flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.Equal(t, "test_key", mockUseCase.recievedKey)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.SuccessLogs)
	assert.Count(t, 1, mockLogger.OutputLogs)
	assert.Contains(t, "test_value", mockLogger.OutputLogs)
}

func TestGetCommand_UseCaseError(t *testing.T) {
	mockUseCase := &mockGetUseCase{
		executeFunc: func(vctx *model.VaultContext, key string) (string, error) {
			return "", fmt.Errorf("%s", test.ErrMsgExecuteFailed)
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewGetCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	cmd.Flags().Bool("cache", false, "Cache passphrase")

	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("key", "test_key"); err != nil {
		t.Fatalf("failed to set key flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, test.ErrMsgExecuteFailed, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.ErrorLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}
