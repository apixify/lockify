package cache

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

const testPassphrase = "test_passphrase"

type mockCachePassphraseUseCase struct {
	executeFunc func(vctx *model.VaultContext, passphrase string) error
}

func (m *mockCachePassphraseUseCase) Execute(vctx *model.VaultContext, passphrase string) error {
	if m.executeFunc != nil {
		return m.executeFunc(vctx, passphrase)
	}
	return nil
}

func TestSetCommand_Success(t *testing.T) {
	mockUseCase := &mockCachePassphraseUseCase{}
	mockPrompt := &test.MockPromptService{
		GetPassphraseInputFunc: func(message string) (string, error) {
			return testPassphrase, nil
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, err := NewSetCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err != nil {
		t.Fatalf("NewSetCommand returned error: %v", err)
	}

	cmd.Flags().Bool("cache", false, "Cache passphrase")
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

func TestSetCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockCachePassphraseUseCase{}
	mockPrompt := &test.MockPromptService{}
	mockLogger := &test.MockLogger{}

	cmd, err := NewSetCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err != nil {
		t.Fatalf("NewSetCommand returned error: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err = cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
	assert.Count(t, 0, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestSetCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockCachePassphraseUseCase{}
	mockPrompt := &test.MockPromptService{}
	mockLogger := &test.MockLogger{}

	cmd, err := NewSetCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err != nil {
		t.Fatalf("NewSetCommand returned error: %v", err)
	}

	if err := cmd.Flags().Set("env", ""); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err = cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
	assert.Count(t, 0, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestSetCommand_Error_PromptFailed(t *testing.T) {
	mockUseCase := &mockCachePassphraseUseCase{}
	mockPrompt := &test.MockPromptService{
		GetPassphraseInputFunc: func(message string) (string, error) {
			return "", fmt.Errorf("prompt failed")
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, err := NewSetCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err != nil {
		t.Fatalf("NewSetCommand returned error: %v", err)
	}

	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err = cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, "prompt failed", err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestSetCommand_Error_EmptyPassphrase(t *testing.T) {
	mockUseCase := &mockCachePassphraseUseCase{}
	mockPrompt := &test.MockPromptService{
		GetPassphraseInputFunc: func(message string) (string, error) {
			return "", nil
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, err := NewSetCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err != nil {
		t.Fatalf("NewSetCommand returned error: %v", err)
	}

	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err = cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, "passphrase cannot be empty", err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestSetCommand_UseCaseError(t *testing.T) {
	expectedError := "cache failed"
	mockUseCase := &mockCachePassphraseUseCase{
		executeFunc: func(vctx *model.VaultContext, passphrase string) error {
			return fmt.Errorf("%s", expectedError)
		},
	}
	mockPrompt := &test.MockPromptService{
		GetPassphraseInputFunc: func(message string) (string, error) {
			return testPassphrase, nil
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, err := NewSetCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err != nil {
		t.Fatalf("NewSetCommand returned error: %v", err)
	}

	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err = cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, expectedError, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.ErrorLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestSetCommand_VerifiesPassphraseMessage(t *testing.T) {
	env := "production"
	expectedMessage := fmt.Sprintf("Enter passphrase for environment %q:", env)

	var receivedMessage string
	mockUseCase := &mockCachePassphraseUseCase{}
	mockPrompt := &test.MockPromptService{
		GetPassphraseInputFunc: func(message string) (string, error) {
			receivedMessage = message
			return testPassphrase, nil
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, err := NewSetCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err != nil {
		t.Fatalf("NewSetCommand returned error: %v", err)
	}

	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", env); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("RunE returned unexpected error: %v", err)
	}

	if receivedMessage != expectedMessage {
		t.Errorf("Prompt message = %q, want %q", receivedMessage, expectedMessage)
	}
}

func TestSetCommand_PassesCorrectDataToUseCase(t *testing.T) {
	expectedEnv := "staging"
	expectedPassphrase := "secret123"

	var receivedEnv, receivedPassphrase string
	mockUseCase := &mockCachePassphraseUseCase{
		executeFunc: func(vctx *model.VaultContext, passphrase string) error {
			receivedEnv = vctx.Env
			receivedPassphrase = passphrase
			return nil
		},
	}
	mockPrompt := &test.MockPromptService{
		GetPassphraseInputFunc: func(message string) (string, error) {
			return expectedPassphrase, nil
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, err := NewSetCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err != nil {
		t.Fatalf("NewSetCommand returned error: %v", err)
	}

	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", expectedEnv); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	if err := cmd.RunE(cmd, nil); err != nil {
		t.Fatalf("RunE returned unexpected error: %v", err)
	}

	if receivedEnv != expectedEnv {
		t.Errorf("UseCase received env = %q, want %q", receivedEnv, expectedEnv)
	}
	if receivedPassphrase != expectedPassphrase {
		t.Errorf(
			"UseCase received passphrase = %q, want %q",
			receivedPassphrase,
			expectedPassphrase,
		)
	}
}
