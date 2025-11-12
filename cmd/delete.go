package cmd

import (
	"errors"
	"fmt"

	"github.com/apixify/lockify/internal/service"
	"github.com/apixify/lockify/internal/vault"
	"github.com/spf13/cobra"
)

// lockify del --env [env] --key [key]
var delCmd = &cobra.Command{
	Use:   "del",
	Short: "delete an entry from the vault of an env",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("⏳ removing key...")
		env, err := cmd.Flags().GetString("env")
		if err != nil {
			return fmt.Errorf("failed to retrieve key flag")
		}
		if env == "" {
			return fmt.Errorf("env is required")
		}

		key, err := cmd.Flags().GetString("key")
		if err != nil {
			return fmt.Errorf("failed to retrieve key flag")
		}
		if key == "" {
			return fmt.Errorf("key is required")
		}

		vault, err := vault.Open(env)
		if err != nil {
			return fmt.Errorf("failed to open vault for environment %s: %w", env, err)
		}

		passphraseService := service.NewPassphraseService(env)
		passphrase := passphraseService.GetPassphrase()
		if !vault.VerifyFingerPrint(passphrase) {
			passphraseService.ClearPassphrase()
			return errors.New("invalid credentials")
		}

		vault.DeleteEntry(key)
		err = vault.Save()
		if err != nil {
			return fmt.Errorf("failed to save vault")
		}

		fmt.Printf("✅ key %s is removed successfully.\n", key)

		return nil
	},
}

func init() {
	delCmd.Flags().StringP("env", "e", "", "Environment Name")
	delCmd.Flags().StringP("key", "k", "", "key to delete from the vault")
	rootCmd.AddCommand(delCmd)
}
