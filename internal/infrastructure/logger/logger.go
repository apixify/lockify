package logger

import (
	"fmt"
	"io"
	"os"
)

// Logger provides structured logging
type Logger struct {
	stderr io.Writer
	stdout io.Writer
}

// New creates a new logger
func New() *Logger {
	return &Logger{
		stderr: os.Stderr,
		stdout: os.Stdout,
	}
}

// Info writes an info message to stderr
func (l *Logger) Info(format string, args ...interface{}) {
	fmt.Fprintf(l.stderr, format+"\n", args...)
}

// Error writes an error message to stderr
func (l *Logger) Error(format string, args ...interface{}) {
	fmt.Fprintf(l.stderr, "❌ "+format+"\n", args...)
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

// OutputRaw writes raw data to stdout
func (l *Logger) OutputRaw(data []byte) {
	l.stdout.Write(data)
}
