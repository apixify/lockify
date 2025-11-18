package domain

type Logger interface {
	// Info writes an info message to logs
	Info(format string, args ...interface{})
	// Error writes an error message to logs
	Error(format string, args ...interface{})
	// Warning writes a warning message to logs
	Warning(format string, args ...interface{})
	// Success writes a success message to logs
	Success(format string, args ...interface{})
	// Progress writes a progress message to logs
	Progress(format string, args ...interface{})
	// Output writes to stdout (for data output, not logs)
	Output(format string, args ...interface{})
}
