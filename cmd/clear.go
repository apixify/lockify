package cmd

import (
	"github.com/apixify/lockify/internal/di"
	"github.com/spf13/cobra"
)

// lockify cache clear
var clearCmd = &cobra.Command{
	Use:   "cache clear",
	Short: "Clear cached passphrases",
	Long: `Clear all cached passphrases from the system keyring.

This command removes all passphrases that were cached in the system keyring.
You will be prompted for passphrases again on next use.`,
	Example: `  lockify cache clear`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Progress("clearing cached passphrases")
		useCase := di.BuildClearCachedPassphrase()

		ctx := getContext()
		err := useCase.Execute(ctx)
		if err != nil {
			logger.Error("failed to cleare cached passphrases")
			return err
		}

		logger.Success("cleared cached passphrases")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
