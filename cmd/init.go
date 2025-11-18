package cmd

import (
	"github.com/apixify/lockify/internal/di"
	"github.com/spf13/cobra"
)

// lockify init --env [env]
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Lockify vault in the current directory",
	Long: `Initialize a new Lockify vault for an environment.

This command creates a new encrypted vault file that will store your environment variables.
You will be prompted for a passphrase that will be used to encrypt and decrypt your secrets.`,
	Example: `  lockify init --env prod
  lockify init --env staging
  lockify init -e local`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := requireEnvFlag(cmd)
		if err != nil {
			return err
		}

		logger.Progress("Initializing Lockify vault")
		ctx := getContext()
		useCase := di.BuildInitializeVault()
		vault, err := useCase.Execute(ctx, env)
		if err != nil {
			return err
		}

		logger.Success("Lockify vault initialized at %s", vault.Path())
		return nil
	},
}

func init() {
	initCmd.Flags().StringP("env", "e", "", "Environment Name")
	initCmd.MarkFlagRequired("env")

	rootCmd.AddCommand(initCmd)
}
