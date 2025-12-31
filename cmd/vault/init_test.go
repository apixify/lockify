package vault

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

type mockInitUseCase struct {
	executeFunc func(ctx context.Context, env string, shouldCache bool) (*model.Vault, error)
}

func (m *mockInitUseCase) Execute(
	ctx context.Context,
	env string,
	shouldCache bool,
) (*model.Vault, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, env, shouldCache)
	}
	vault, _ := model.NewVault(env, "finger", "salt")
	vault.SetPath("/tmp/test.vault")
	return vault, nil
}

func TestInitCommand_Success(t *testing.T) {
	mockUseCase := &mockInitUseCase{}
	mockLogger := &test.MockLogger{}
	cmdCtx := cli.NewCommandContext()

	cmd, _ := NewInitCommand(mockUseCase, mockLogger, cmdCtx)
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("RunE returned unexpected error: %v", err)
	}

	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.SuccessLogs)
}

func TestInitCommand_Failed(t *testing.T) {
	errMsg := "error during execution"
	mockUseCase := &mockInitUseCase{
		executeFunc: func(ctx context.Context, env string, shouldCache bool) (*model.Vault, error) {
			return nil, fmt.Errorf("%s", errMsg)
		},
	}
	mockLogger := &test.MockLogger{}
	cmdCtx := cli.NewCommandContext()

	cmd, _ := NewInitCommand(mockUseCase, mockLogger, cmdCtx)
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
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

func TestInitCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockInitUseCase{}
	mockLogger := &test.MockLogger{}
	cmdCtx := cli.NewCommandContext()

	cmd, _ := NewInitCommand(mockUseCase, mockLogger, cmdCtx)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
	assert.Count(t, 0, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestInitCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockInitUseCase{}
	mockLogger := &test.MockLogger{}
	cmdCtx := cli.NewCommandContext()

	cmd, _ := NewInitCommand(mockUseCase, mockLogger, cmdCtx)
	if err := cmd.Flags().Set("env", ""); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
	assert.Count(t, 0, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}
