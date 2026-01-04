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

type mockDeleteUseCase struct {
	executeFunc func(vctx *model.VaultContext, key string) error
	receivedEnv string
	receivedKey string
}

func (m *mockDeleteUseCase) Execute(vctx *model.VaultContext, key string) error {
	m.receivedEnv = vctx.Env
	m.receivedKey = key
	if m.executeFunc != nil {
		return m.executeFunc(vctx, key)
	}
	return nil
}

func TestDeleteCommand_Success(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewDeleteCommand(mockUseCase, mockLogger, cli.NewCommandContext())
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
	assert.Equal(t, "test", mockUseCase.receivedEnv)
	assert.Equal(t, "test_key", mockUseCase.receivedKey)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.SuccessLogs)
	assert.Contains(t, "test_key", mockLogger.SuccessLogs[0])
}

func TestDeleteCommand_UseCaseError(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{
		executeFunc: func(vctx *model.VaultContext, key string) error {
			return fmt.Errorf("%s", test.ErrMsgExecuteFailed)
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewDeleteCommand(mockUseCase, mockLogger, cli.NewCommandContext())
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
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestDeleteCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewDeleteCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("key", "test_key"); err != nil {
		t.Fatalf("failed to set key flag: %v", err)
	}

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

func TestDeleteCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewDeleteCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("env", ""); err != nil {
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
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
	// Progress is logged before flag validation
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestDeleteCommand_Error_Required_Key(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewDeleteCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	wants := "key flag is required"
	assert.Contains(t, wants, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestDeleteCommand_Error_Empty_Key(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewDeleteCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("key", ""); err != nil {
		t.Fatalf("failed to set key flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	wants := "key flag is required"
	assert.Contains(t, wants, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}
