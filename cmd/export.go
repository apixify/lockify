package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/apixify/lockify/internal/service"
	"github.com/apixify/lockify/internal/vault"
	"github.com/spf13/cobra"
)

// lockify export --env [env] --format [dotenv|json]
// lockify export --env prod --format dotenv > .env
// lockify export --env staging --format json > env.json
const (
	dotenvFormat = "dotenv"
	jsonFormat   = "json"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Decrypt all variables and export them in a specific format.",
	RunE: func(cmd *cobra.Command, args []string) error {
		envPassphraseKey, _ := cmd.Flags().GetString("passphrase-env")
		env, err := cmd.Flags().GetString("env")
		if err != nil {
			return fmt.Errorf("failed to retrieve env flag")
		}
		if env == "" {
			return fmt.Errorf("env is required")
		}
		fmt.Fprintf(os.Stderr, "⏳ Exporting secrets for %s...\n", env)

		format, err := cmd.Flags().GetString("format")
		if err != nil {
			return fmt.Errorf("failed to retrieve format flag")
		}
		if format == "" {
			return fmt.Errorf("format is required")
		}
		if format != dotenvFormat && format != jsonFormat {
			return fmt.Errorf("format must be either %s or %s. %s is given", dotenvFormat, jsonFormat, format)
		}

		passphraseService := service.NewPassphraseService(env, envPassphraseKey)
		vault, err := vault.Open(env)
		if err != nil {
			return fmt.Errorf("failed to open vault for environment %s: %w", env, err)
		}

		passphrase := passphraseService.GetPassphrase()
		if !vault.VerifyFingerPrint(passphrase) {
			passphraseService.ClearPassphrase()
			return fmt.Errorf("invalid credentials")
		}

		crypto, err := service.NewCryptoService(vault.Meta.Salt, passphrase)
		if err != nil {
			return fmt.Errorf("failed to initialize crypto service: %w", err)
		}

		if format == dotenvFormat {
			for k, v := range vault.Entries {
				decryptedVal, _ := crypto.DecryptValue(v.Value)
				fmt.Fprintf(os.Stdout, "%s=%s\n", k, decryptedVal)
			}
		} else {
			mappedEntries := make(map[string]string)
			for k, v := range vault.Entries {
				decryptedVal, _ := crypto.DecryptValue(v.Value)
				mappedEntries[k] = decryptedVal
			}

			data, _ := json.MarshalIndent(mappedEntries, "", "  ")
			fmt.Println(string(data))
		}

		return nil
	},
}

func init() {
	exportCmd.Flags().StringP("env", "e", "", "Environment Name")
	exportCmd.Flags().String("format", "dotenv", "The format of the exported file [dotenv|json]")
	exportCmd.Flags().String("passphrase-env", "LOCKIFY_PASSPHRASE", "Name of the environment variable that holds the passphrase")

	rootCmd.AddCommand(exportCmd)
}
