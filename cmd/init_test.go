package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/test"
)

type mockInitUseCase struct {
	executeFunc func(ctx context.Context, env string) (*model.Vault, error)
}

func (m *mockInitUseCase) Execute(ctx context.Context, env string) (*model.Vault, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, env)
	}
	vault, _ := model.NewVault(env, "finger", "salt")
	vault.SetPath("/tmp/test.vault")
	return vault, nil
}

func TestInitCommand_Success(t *testing.T) {
	mockUseCase := &mockInitUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewInitCommand(mockUseCase, mockLogger)
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("RunE returned unexpected error: %v", err)
	}

	if len(mockLogger.ProgressLogs) == 0 {
		t.Error("expected Progress to be logged")
	}
	if len(mockLogger.SuccessLogs) == 0 {
		t.Error("expected Success to be logged")
	}
}

func TestInitCommand_Failed(t *testing.T) {
	errMsg := "error during execution"
	mockUseCase := &mockInitUseCase{executeFunc: func(ctx context.Context, env string) (*model.Vault, error) {
		return nil, fmt.Errorf("%s", errMsg)
	}}
	mockLogger := &test.MockLogger{}

	cmd := NewInitCommand(mockUseCase, mockLogger)
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("RunE expected error")
	}

	if !strings.Contains(err.Error(), errMsg) {
		t.Errorf("wants %q, got %q", errMsg, err)
	}
	if len(mockLogger.ProgressLogs) == 0 {
		t.Error("expected Progress to be logged")
	}
	if len(mockLogger.SuccessLogs) != 0 {
		t.Error("is not expecting Success to be logged")
	}
}

func TestInitCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockInitUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewInitCommand(mockUseCase, mockLogger)

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatalf("RunE is expecting error: %v", err)
	}

	wants := "env flag is required (use --env or -e)"
	if !strings.Contains(err.Error(), wants) {
		t.Errorf("wants error to contain %q, got %q", wants, err)
	}
	if len(mockLogger.ProgressLogs) != 0 {
		t.Error("is not expecting Progress to be logged")
	}
	if len(mockLogger.SuccessLogs) != 0 {
		t.Error("is not expecting Success to be logged")
	}
}

func TestInitCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockInitUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewInitCommand(mockUseCase, mockLogger)
	if err := cmd.Flags().Set("env", ""); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatalf("RunE is expecting error: %v", err)
	}

	wants := "env flag is required (use --env or -e)"
	if !strings.Contains(err.Error(), wants) {
		t.Errorf("wants error to contain %q, got %q", wants, err)
	}
	if len(mockLogger.ProgressLogs) != 0 {
		t.Error("is not expecting Progress to be logged")
	}
	if len(mockLogger.SuccessLogs) != 0 {
		t.Error("is not expecting Success to be logged")
	}
}
