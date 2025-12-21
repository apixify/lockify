package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

type mockImportUseCase struct {
	executeFunc       func(ctx context.Context, env string, format value.FileFormat, r io.Reader, overwrite bool) (int, int, error)
	receivedEnv       string
	receivedFormat    value.FileFormat
	receivedOverwrite bool
}

func (m *mockImportUseCase) Execute(
	ctx context.Context,
	env string,
	format value.FileFormat,
	r io.Reader,
	overwrite bool,
) (int, int, error) {
	m.receivedEnv = env
	m.receivedFormat = format
	m.receivedOverwrite = overwrite
	if m.executeFunc != nil {
		return m.executeFunc(ctx, env, format, r, overwrite)
	}
	return 3, 1, nil
}

func TestImportCommand_Success(t *testing.T) {
	mockUseCase := &mockImportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewImportCommand(mockUseCase, mockLogger)
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("format", "dotenv"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	// Simulate stdin by not providing file args
	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.Equal(t, "test", mockUseCase.receivedEnv)
	assert.Equal(t, value.DotEnv, mockUseCase.receivedFormat)
	assert.False(t, mockUseCase.receivedOverwrite)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.SuccessLogs)
	assert.Contains(t, "Imported 3 key(s), skipped 1 key(s)", mockLogger.SuccessLogs[0])
}

func TestImportCommand_Success_WithOverwrite(t *testing.T) {
	mockUseCase := &mockImportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewImportCommand(mockUseCase, mockLogger)
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("format", "json"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}
	if err := cmd.Flags().Set("overwrite", "true"); err != nil {
		t.Fatalf("failed to set overwrite flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
	assert.Equal(t, "test", mockUseCase.receivedEnv)
	assert.Equal(t, value.JSON, mockUseCase.receivedFormat)
	assert.True(t, mockUseCase.receivedOverwrite)
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 1, mockLogger.SuccessLogs)
}

func TestImportCommand_UseCaseError(t *testing.T) {
	mockUseCase := &mockImportUseCase{
		executeFunc: func(ctx context.Context, env string, format value.FileFormat, r io.Reader, overwrite bool) (int, int, error) {
			return 0, 0, fmt.Errorf("%s", errMsgExecuteFailed)
		},
	}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewImportCommand(mockUseCase, mockLogger)
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
	assert.Contains(t, errMsgExecuteFailed, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
	assert.Count(t, 0, mockLogger.SuccessLogs)
}

func TestImportCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockImportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewImportCommand(mockUseCase, mockLogger)
	if err := cmd.Flags().Set("format", "dotenv"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	assert.Contains(t, errMsgEmptyEnv, err.Error())
}

func TestImportCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockImportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewImportCommand(mockUseCase, mockLogger)
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
	assert.Contains(t, errMsgEmptyEnv, err.Error())
}

func TestImportCommand_Error_Empty_Format(t *testing.T) {
	mockUseCase := &mockImportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewImportCommand(mockUseCase, mockLogger)
	if err := cmd.Flags().Set("env", "test"); err != nil {
		t.Fatalf("failed to set env flag: %v", err)
	}
	if err := cmd.Flags().Set("format", ""); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	wants := "format flag is required"
	assert.Contains(t, wants, err.Error())
}

func TestImportCommand_Error_Invalid_Format(t *testing.T) {
	mockUseCase := &mockImportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd, _ := NewImportCommand(mockUseCase, mockLogger)
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
