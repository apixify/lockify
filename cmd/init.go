package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/apixify/lockify/internal/service"
	"github.com/apixify/lockify/internal/vault"
	"github.com/spf13/cobra"
)

// lockify init --env [env]
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Lockify vault in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := cmd.Flags().GetString("env")
		if err != nil {
			return fmt.Errorf("failed to retrieve key flag")
		}
		if env == "" {
			return fmt.Errorf("env is required")
		}

		vaultPath := filepath.Join(".lockify", env+".vault.enc")

		fmt.Println("⏳ Initializing Lockify vault at", vaultPath)
		if _, err := os.Stat(vaultPath); err == nil {
			return fmt.Errorf("vault already exists at %s", vaultPath)
		}

		if err := os.MkdirAll(".lockify", 0700); err != nil {
			return fmt.Errorf("failed to create .lockify directory: %w", err)
		}

		fmt.Println("Creating empty encrypted vault placeholder at", vaultPath)
		passphrase := service.NewPassphraseService(env)
		salt, err := service.GenerateSalt(16)
		if err != nil {
			return fmt.Errorf("failed to generate salt")
		}

		_, err = vault.Create(vaultPath, env, passphrase.GetPassphrase(), salt)
		if err != nil {
			return fmt.Errorf("failed to create %s: %w", vaultPath, err)
		}

		fmt.Println("✅ Lockify vault initialized at .lockify/")
		return nil
	},
}

func init() {
	initCmd.Flags().StringP("env", "e", "", "Environment Name")
	rootCmd.AddCommand(initCmd)
}
