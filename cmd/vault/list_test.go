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

type mockListUseCase struct {
	executeFunc func(vctx *model.VaultContext) ([]string, error)
	receivedEnv string
}

func (m *mockListUseCase) Execute(vctx *model.VaultContext) ([]string, error) {
	m.receivedEnv = vctx.Env
	if m.executeFunc != nil {
		return m.executeFunc(vctx)
	}
	return []string{"key1", "key2", "key3"}, nil
}

func TestListCommand_Success(t *testing.T) {
	mockUseCase := &mockListUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewListCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.Equal(t, "test", mockUseCase.receivedEnv)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.SuccessLogs)
	assert.Count(t, 3, mockLogger.OutputLogs)
	assert.Contains(t, "  - key1\n", mockLogger.OutputLogs)
	assert.Contains(t, "  - key2\n", mockLogger.OutputLogs)
	assert.Contains(t, "  - key3\n", mockLogger.OutputLogs)
}

func TestListCommand_EmptyKeys(t *testing.T) {
	mockUseCase := &mockListUseCase{
		executeFunc: func(vctx *model.VaultContext) ([]string, error) {
			return []string{}, nil
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewListCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
	assert.Count(t, 1, mockLogger.InfoLogs)
	assert.Contains(t, "No entries found", mockLogger.InfoLogs[0])
}

func TestListCommand_UseCaseError(t *testing.T) {
	mockUseCase := &mockListUseCase{
		executeFunc: func(vctx *model.VaultContext) ([]string, error) {
			return nil, fmt.Errorf("%s", test.ErrMsgExecuteFailed)
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewListCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, test.ErrMsgExecuteFailed, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestListCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockListUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewListCommand(mockUseCase, mockLogger, cli.NewCommandContext())

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
	// Progress is logged before flag validation
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestListCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockListUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewListCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("env", ""); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}
