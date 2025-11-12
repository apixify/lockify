package cmd

import (
	"errors"
	"fmt"

	"github.com/apixify/lockify/internal/service"
	"github.com/apixify/lockify/internal/vault"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [env]",
	Short: "get a secret from the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("getting a secret from the vault")
		env := args[0]
		key, _ := cmd.Flags().GetString("key")
		passphraseService := service.NewPassphraseService(env)
		vault, err := vault.Open(env)
		if err != nil {
			return fmt.Errorf("failed to open vault for environment %s: %w", env, err)
		}

		passphrase := passphraseService.GetPassphrase()
		if !vault.VerifyFingerPrint(passphrase) {
			passphraseService.ClearPassphrase()
			return errors.New("invalid credentials")
		}

		entry, err := vault.GetEntry(key)
		if err != nil {
			return fmt.Errorf("%w", err)
		}

		crypto, err := service.NewCryptoService(vault.Meta.Salt, passphrase)
		if err != nil {
			return fmt.Errorf("failed to initialize the crypto service")
		}

		value, _ := crypto.DecryptValue(entry.Value)

		fmt.Println(value)

		return nil
	},
}

func init() {
	getCmd.Flags().StringP("key", "k", "", "The key to use for getting the secret")

	rootCmd.AddCommand(getCmd)
}
