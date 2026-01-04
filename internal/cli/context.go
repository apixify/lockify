package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// CommandContext provides shared utilities for command execution
type CommandContext struct{}

// NewCommandContext creates a new CommandContext with default implementations
func NewCommandContext() *CommandContext {
	return &CommandContext{}
}

// GetContext returns a context for command execution (unexported)
func (c *CommandContext) GetContext() context.Context {
	return context.Background()
}

// RequireEnvFlag retrieves the env flag and returns an error if empty (unexported)
func (c *CommandContext) RequireEnvFlag(cmd *cobra.Command) (string, error) {
	env, err := cmd.Flags().GetString("env")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve env flag: %w", err)
	}
	if env == "" {
		return "", errors.New(ErrMsgEmptyEnv)
	}
	return env, nil
}

// RequireStringFlag retrieves a string flag and returns an error if empty (unexported)
func (c *CommandContext) RequireStringFlag(cmd *cobra.Command, flag string) (string, error) {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve %s flag: %w", flag, err)
	}
	if value == "" {
		return "", fmt.Errorf("%s flag is required", flag)
	}
	return value, nil
}

// GetCacheFlag retrieves the --cache flag value (unexported)
func (c *CommandContext) GetCacheFlag(cmd *cobra.Command) (bool, error) {
	cache, err := cmd.Flags().GetBool("cache")
	if err != nil {
		return false, fmt.Errorf("failed to retrieve cache flag: %w", err)
	}
	return cache, nil
}
