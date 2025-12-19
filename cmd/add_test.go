package cmd

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/test"
)

type mockAddUseCase struct {
	executeFunc func(ctx context.Context, dto app.AddEntryDTO) error
	receivedDTO app.AddEntryDTO
}

func (m *mockAddUseCase) Execute(ctx context.Context, dto app.AddEntryDTO) error {
	m.receivedDTO = dto
	if m.executeFunc != nil {
		return m.executeFunc(ctx, dto)
	}

	return nil
}

func TestAddCommand_Success(t *testing.T) {
	mockUseCase := &mockAddUseCase{}
	mockLogger := &test.MockLogger{}
	mockPrompt := &test.MockPromptService{
		GetUserInputFunc: func(isSecret bool) (string, string) {
			return "test_key", "test_value"
		},
	}

	cmd := NewAddCommand(mockUseCase, mockPrompt, mockLogger)
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("RunE returned unexpected error: %v", err)
	}

	if mockUseCase.receivedDTO.Env != "test" || mockUseCase.receivedDTO.Key != "test_key" || mockUseCase.receivedDTO.Value != "test_value" {
		t.Fatalf("unexpected DTO: %+v", mockUseCase.receivedDTO)
	}
	if len(mockLogger.ProgressLogs) == 0 {
		t.Error("expected Progress to be logged")
	}
	if len(mockLogger.SuccessLogs) == 0 {
		t.Error("expected Success to be logged")
	}
}

func TestAddCommand_UseCaseError(t *testing.T) {
	errMsg := "execute failed"
	mockUseCase := &mockAddUseCase{
		executeFunc: func(ctx context.Context, dto app.AddEntryDTO) error {
			return fmt.Errorf("%s", errMsg)
		},
	}
	mockLogger := &test.MockLogger{}
	mockPrompt := &test.MockPromptService{
		GetUserInputFunc: func(isSecret bool) (string, string) {
			return "test_key", "test_value"
		},
	}

	cmd := NewAddCommand(mockUseCase, mockPrompt, mockLogger)

	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), errMsg) {
		t.Fatalf("expected error to contain %q, got %v", errMsg, err)
	}
	if len(mockLogger.ProgressLogs) == 0 {
		t.Error("expected Progress to be logged")
	}
	if len(mockLogger.SuccessLogs) != 0 {
		t.Error("did not expect Success to be logged")
	}
	if len(mockLogger.ErrorLogs) == 0 {
		t.Error("expected Error to be logged")
	}
}
