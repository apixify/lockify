package cmd

import (
	"fmt"
	"os"

	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/spf13/cobra"
)

// NewRootCmd returns the root command for subcommand registration
func NewRootCmd() (*cobra.Command, *cli.CommandContext) {
	cmdCtx := cli.NewCommandContext()
	rootCmd := &cobra.Command{
		Use:   "lockify",
		Short: "Lockify securely manages your .env files and secrets",
		Long: `Lockify is a lightweight CLI tool for securely managing environment variables and .env files locally.
	
		Lockify encrypts your environment variables using AES-GCM encryption with Argon2 key derivation.
		Your secrets are protected with a passphrase that can be stored securely in your system's keyring.`,
		Version: Version,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(os.Stderr, "Welcome to Lockify! Use --help to see available commands.")
		},
	}
	// Add global --cache flag
	rootCmd.PersistentFlags().BoolP("cache", "c", false, "Cache passphrase in system keyring")

	// Register built-in commands
	versionCmd := NewVersionCommand(di.GetLogger())
	rootCmd.AddCommand(versionCmd)

	return rootCmd, cmdCtx
}
