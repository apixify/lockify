package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/apixify/lockify/internal/domain"
)

// Logger provides structured logging
type Logger struct {
	stderr io.Writer
	stdout io.Writer
}

// New creates a new logger
func New() domain.Logger {
	return &Logger{
		stderr: os.Stderr,
		stdout: os.Stdout,
	}
}

// Info writes an info message to stderr
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Fprintf(l.stderr, "ℹ️"+format+"\n", args...)
}

// Error writes an error message to stderr
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Fprintf(l.stderr, "❌ "+format+"\n", args...)
}

// Warning writes a warning message to stderr
func (l *Logger) Warning(format string, args ...interface{}) {
	fmt.Fprintf(l.stderr, "🔶 "+format+"\n", args...)
}

// Success writes a success message to stderr
func (l *Logger) Success(format string, args ...interface{}) {
	fmt.Fprintf(l.stderr, "✅ "+format+"\n", args...)
}

// Progress writes a progress message to stderr
func (l *Logger) Progress(format string, args ...interface{}) {
	fmt.Fprintf(l.stderr, "⏳ "+format+"\n", args...)
}

// Output writes to stdout (for data output, not logs)
func (l *Logger) Output(format string, args ...interface{}) {
	fmt.Fprintf(l.stdout, format+"\n", args...)
}
