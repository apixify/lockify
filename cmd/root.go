package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "lockify",
	Short: "Lockify securely manages your .env files and secrets",
	Long: `Lockify is a lightweight CLI tool for securely managing environment variables and .env files locally.

Lockify encrypts your environment variables using AES-GCM encryption with Argon2 key derivation.
Your secrets are protected with a passphrase that can be stored securely in your system's keyring.`,
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(os.Stderr, "Welcome to Lockify! Use --help to see available commands.")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

// requireEnvFlag retrieves the env flag and returns an error if empty
func requireEnvFlag(cmd *cobra.Command) (string, error) {
	env, err := cmd.Flags().GetString("env")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve env flag: %w", err)
	}
	if env == "" {
		return "", fmt.Errorf("env flag is required (use --env or -e)")
	}
	return env, nil
}

func requireStringFlag(cmd *cobra.Command, flag string) (string, error) {
	value, err := cmd.Flags().GetString(flag)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve %s flag: %w", flag, err)
	}
	if value == "" {
		return "", fmt.Errorf("%s flag is required", flag)
	}
	return value, nil
}

// getContext returns a context for command execution
func getContext() context.Context {
	return context.Background()
}
