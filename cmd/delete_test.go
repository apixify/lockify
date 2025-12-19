package cmd

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

type mockDeleteUseCase struct {
	executeFunc func(ctx context.Context, env, key string) error
	receivedEnv string
	receivedKey string
}

func (m *mockDeleteUseCase) Execute(ctx context.Context, env, key string) error {
	m.receivedEnv = env
	m.receivedKey = key
	if m.executeFunc != nil {
		return m.executeFunc(ctx, env, key)
	}
	return nil
}

func TestDeleteCommand_Success(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewDeleteCommand(mockUseCase, mockLogger)
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
	errMsg := "execute failed"
	mockUseCase := &mockDeleteUseCase{
		executeFunc: func(ctx context.Context, env, key string) error {
			return fmt.Errorf("%s", errMsg)
		},
	}
	mockLogger := &test.MockLogger{}

	cmd := NewDeleteCommand(mockUseCase, mockLogger)
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
	assert.Contains(t, errMsg, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestDeleteCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewDeleteCommand(mockUseCase, mockLogger)
	if err := cmd.Flags().Set("key", "test_key"); err != nil {
		t.Fatalf("failed to set key flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	wants := "env flag is required (use --env or -e)"
	assert.Contains(t, wants, err.Error())
	// Progress is logged before flag validation
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestDeleteCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewDeleteCommand(mockUseCase, mockLogger)
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
	wants := "env flag is required (use --env or -e)"
	assert.Contains(t, wants, err.Error())
	// Progress is logged before flag validation
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestDeleteCommand_Error_Required_Key(t *testing.T) {
	mockUseCase := &mockDeleteUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewDeleteCommand(mockUseCase, mockLogger)
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

	cmd := NewDeleteCommand(mockUseCase, mockLogger)
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
