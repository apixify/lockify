package vault

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

type mockExportUseCase struct {
	executeFunc    func(vctx *model.VaultContext, exportFormat value.FileFormat) error
	receivedEnv    string
	receivedFormat value.FileFormat
}

func (m *mockExportUseCase) Execute(
	vctx *model.VaultContext,
	exportFormat value.FileFormat,
) error {
	m.receivedEnv = vctx.Env
	m.receivedFormat = exportFormat
	if m.executeFunc != nil {
		return m.executeFunc(vctx, exportFormat)
	}
	return nil
}

func TestExportCommand_Success_DotEnv(t *testing.T) {
	mockUseCase := &mockExportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewExportCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("format", "dotenv"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.Equal(t, "test", mockUseCase.receivedEnv)
	assert.Equal(t, value.DotEnv, mockUseCase.receivedFormat)
	assert.Count(t, 1, mockLogger.ProgressLogs)
}

func TestExportCommand_Success_Json(t *testing.T) {
	mockUseCase := &mockExportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewExportCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("format", "json"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.Equal(t, "test", mockUseCase.receivedEnv)
	assert.Equal(t, value.JSON, mockUseCase.receivedFormat)
	assert.Count(t, 1, mockLogger.ProgressLogs)
}

func TestExportCommand_UseCaseError(t *testing.T) {
	mockUseCase := &mockExportUseCase{
		executeFunc: func(vctx *model.VaultContext, exportFormat value.FileFormat) error {
			return fmt.Errorf("%s", test.ErrMsgExecuteFailed)
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewExportCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	cmd.Flags().Bool("cache", false, "Cache passphrase")
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("format", "dotenv"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, test.ErrMsgExecuteFailed, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
}

func TestExportCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockExportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewExportCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("format", "dotenv"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
}

func TestExportCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockExportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewExportCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("env", ""); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("format", "dotenv"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, cli.ErrMsgEmptyEnv, err.Error())
}

func TestExportCommand_Error_Invalid_Format(t *testing.T) {
	mockUseCase := &mockExportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewExportCommand(mockUseCase, mockLogger, cli.NewCommandContext())
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("format", "invalid"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, "invalid file format", err.Error())
}
