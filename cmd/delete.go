package cmd

import (
	"github.com/apixify/lockify/internal/di"
	"github.com/spf13/cobra"
)

// lockify del --env [env] --key [key]
var delCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"del", "rm"},
	Short:   "Delete an entry from the vault",
	Long: `Delete an entry from the vault.

This command removes a key-value pair from the vault for the specified environment.`,
	Example: `  lockify delete --env prod --key OLD_KEY
  lockify del --env staging -k DEPRECATED_KEY`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Progress("removing key...")
		env, err := requireEnvFlag(cmd)
		if err != nil {
			return err
		}
		key, err := requireStringFlag(cmd, "key")
		if err != nil {
			return err
		}

		ctx := getContext()
		useCase := di.BuildDeleteEntry()
		err = useCase.Execute(ctx, env, key)
		if err != nil {
			return err
		}

		logger.Success("key %s is removed successfully.\n", key)

		return nil
	},
}

func init() {
	delCmd.Flags().StringP("env", "e", "", "Environment Name")
	delCmd.Flags().StringP("key", "k", "", "key to delete from the vault")
	delCmd.MarkFlagRequired("env")
	delCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(delCmd)
}
