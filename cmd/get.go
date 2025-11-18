package cmd

import (
	"github.com/apixify/lockify/internal/di"
	"github.com/spf13/cobra"
)

// lockify get --env [env] --key [key]
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a decrypted value from the vault",
	Long: `Get a decrypted value from the vault.

This command retrieves and decrypts a value from the vault for the specified key.
The decrypted value is printed to stdout, making it suitable for shell scripting.`,
	Example: `  lockify get --env prod --key DATABASE_URL
  lockify get --env staging -k API_KEY`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Progress("getting an entry from the vault")
		env, err := requireEnvFlag(cmd)
		if err != nil {
			return err
		}

		key, err := requireStringFlag(cmd, "key")
		if err != nil {
			return err
		}

		ctx := getContext()
		useCase := di.BuildGetEntry()
		value, err := useCase.Execute(ctx, env, key)
		if err != nil {
			return err
		}

		logger.Success("retrieved key's value successfully")
		logger.Output(value)

		return nil
	},
}

func init() {
	getCmd.Flags().StringP("env", "e", "", "Environment name")
	getCmd.Flags().StringP("key", "k", "", "The key to use for getting the entry")
	getCmd.MarkFlagRequired("env")

	rootCmd.AddCommand(getCmd)
}
