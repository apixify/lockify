package cmd

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/ahmed-abdelgawad92/lockify/test"
	"github.com/ahmed-abdelgawad92/lockify/test/assert"
)

type mockExportUseCase struct {
	executeFunc    func(ctx context.Context, env string, exportFormat value.FileFormat) error
	receivedEnv    string
	receivedFormat value.FileFormat
}

func (m *mockExportUseCase) Execute(ctx context.Context, env string, exportFormat value.FileFormat) error {
	m.receivedEnv = env
	m.receivedFormat = exportFormat
	if m.executeFunc != nil {
		return m.executeFunc(ctx, env, exportFormat)
	}
	return nil
}

func TestExportCommand_Success_DotEnv(t *testing.T) {
	mockUseCase := &mockExportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewExportCommand(mockUseCase, mockLogger)
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

	cmd := NewExportCommand(mockUseCase, mockLogger)
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
	assert.Equal(t, value.Json, mockUseCase.receivedFormat)
	assert.Count(t, 1, mockLogger.ProgressLogs)
}

func TestExportCommand_UseCaseError(t *testing.T) {
	errMsg := "execute failed"
	mockUseCase := &mockExportUseCase{
		executeFunc: func(ctx context.Context, env string, exportFormat value.FileFormat) error {
			return fmt.Errorf("%s", errMsg)
		},
	}
	mockLogger := &test.MockLogger{}

	cmd := NewExportCommand(mockUseCase, mockLogger)
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
	assert.Contains(t, errMsg, err.Error())
	assert.Count(t, 1, mockLogger.ProgressLogs)
}

func TestExportCommand_Error_Required_Env(t *testing.T) {
	mockUseCase := &mockExportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewExportCommand(mockUseCase, mockLogger)
	if err := cmd.Flags().Set("format", "dotenv"); err != nil {
		t.Fatalf("failed to set format flag: %v", err)
	}

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)

	err := cmd.RunE(cmd, nil)
	assert.NotNil(t, err)
	wants := "env flag is required (use --env or -e)"
	assert.Contains(t, wants, err.Error())
}

func TestExportCommand_Error_Empty_Env(t *testing.T) {
	mockUseCase := &mockExportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewExportCommand(mockUseCase, mockLogger)
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
	wants := "env flag is required (use --env or -e)"
	assert.Contains(t, wants, err.Error())
}

func TestExportCommand_Error_Invalid_Format(t *testing.T) {
	mockUseCase := &mockExportUseCase{}
	mockLogger := &test.MockLogger{}

	cmd := NewExportCommand(mockUseCase, mockLogger)
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
