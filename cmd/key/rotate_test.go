package key

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

type mockRotateUseCase struct {
	executeFunc               func(vctx *model.VaultContext, currentPassphrase, newPassphrase string) error
	receivedEnv               string
	receivedCurrentPassphrase string
	receivedNewPassphrase     string
}

func (m *mockRotateUseCase) Execute(
	vctx *model.VaultContext,
	currentPassphrase, newPassphrase string,
) error {
	m.receivedEnv = vctx.Env
	m.receivedCurrentPassphrase = currentPassphrase
	m.receivedNewPassphrase = newPassphrase
	if m.executeFunc != nil {
		return m.executeFunc(vctx, currentPassphrase, newPassphrase)
	}
	return nil
}

func TestRotateCommand_Success(t *testing.T) {
	mockUseCase := &mockRotateUseCase{}
	mockLogger := &test.MockLogger{}
	mockPrompt := &test.MockPromptService{
		GetPassphraseInputFunc: func(message string) (string, error) {
			if message == "Enter current passphrase:" {
				return "current_pass", nil
			}
			return "new_pass", nil
		},
	}

	cmd, _ := NewRotateCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
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
	assert.Equal(t, "current_pass", mockUseCase.receivedCurrentPassphrase)
	assert.Equal(t, "new_pass", mockUseCase.receivedNewPassphrase)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.SuccessLogs)
}

func TestRotateCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockRotateUseCase{}
	mockLogger := &test.MockLogger{}
	mockPrompt := &test.MockPromptService{}

	cmd, _ := NewRotateCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
}

func TestRotateCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockRotateUseCase{}
	mockLogger := &test.MockLogger{}
	mockPrompt := &test.MockPromptService{}

	cmd, _ := NewRotateCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("env", ""); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
}

func TestRotateCommand_UseCaseError(t *testing.T) {
	mockUseCase := &mockRotateUseCase{
		executeFunc: func(vctx *model.VaultContext, currentPassphrase, newPassphrase string) error {
			return fmt.Errorf("%s", test.ErrMsgExecuteFailed)
		},
	}
	mockLogger := &test.MockLogger{}
	mockPrompt := &test.MockPromptService{
		GetPassphraseInputFunc: func(message string) (string, error) {
			if message == "Enter current passphrase:" {
				return "current_pass", nil
			}
			return "new_pass", nil
		},
	}

	cmd, _ := NewRotateCommand(mockUseCase, mockPrompt, mockLogger, cli.NewCommandContext())
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
