package vault

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/test"
)

var addTestConstants = struct {
	env   string
	key   string
	value string
}{
	env:   "test",
	key:   "test_key",
	value: "test_value",
}

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
		GetUserInputFunc: func(isSecret bool) (string, string, error) {
			return addTestConstants.key, addTestConstants.value, nil
		},
	}

	cmd, _ := NewAddCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("env", addTestConstants.env); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("RunE returned unexpected error: %v", err)
	}

	if mockUseCase.receivedDTO.Env != addTestConstants.env ||
		mockUseCase.receivedDTO.Key != addTestConstants.key ||
		mockUseCase.receivedDTO.Value != addTestConstants.value {
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
	mockUseCase := &mockAddUseCase{
		executeFunc: func(ctx context.Context, dto app.AddEntryDTO) error {
			return fmt.Errorf("%s", test.ErrMsgExecuteFailed)
		},
	}
	mockLogger := &test.MockLogger{}
	mockPrompt := &test.MockPromptService{
		GetUserInputFunc: func(isSecret bool) (string, string, error) {
			return addTestConstants.key, addTestConstants.value, nil
		},
	}

	cmd, _ := NewAddCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())

	if err := cmd.Flags().Set("env", addTestConstants.env); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
	if !strings.Contains(err.Error(), test.ErrMsgExecuteFailed) {
		t.Fatalf("expected error to contain %q, got %v", test.ErrMsgExecuteFailed, err)
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
