package cmd

import (
	"github.com/apixify/lockify/internal/di"
	"github.com/spf13/cobra"
)

// lockify list [env]
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all keys in the vault",
	Long: `List all keys in the vault.

This command displays all keys stored in the vault for the specified environment.
Only keys are displayed, not decrypted values, for security reasons.`,
	Example: `  lockify list --env prod
  lockify list --env staging`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Progress("Listing all secrets in the vault")
		env, err := requireEnvFlag(cmd)
		if err != nil {
			return err
		}

		ctx := getContext()
		useCase := di.BuildListEntries()
		keys, err := useCase.Execute(ctx, env)
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			logger.Info("No entries found in vault")
			return nil
		}

		logger.Success("Found %d key(s):", len(keys))
		for _, v := range keys {
			logger.Output("  - %s\n", v)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().StringP("env", "e", "", "Environment Name")
	listCmd.MarkFlagRequired("env")

	rootCmd.AddCommand(listCmd)
}
