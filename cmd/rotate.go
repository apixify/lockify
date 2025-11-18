package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/apixify/lockify/internal/di"
	"github.com/spf13/cobra"
)

// lockify rotate-key --env [env]
var rotateCmd = &cobra.Command{
	Use:   "rotate-key",
	Short: "Rotate the passphrase for a vault",
	Long: `Rotate the passphrase for a vault.

This command allows you to change the passphrase for a vault by re-encrypting all entries
with a new passphrase. You will be prompted for the current passphrase and a new passphrase.`,
	Example: `  lockify rotate-key --env prod
  lockify rotate-key --env staging`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := requireEnvFlag(cmd)
		if err != nil {
			return err
		}

		var passphrase string
		prompt := &survey.Password{Message: "Enter current passphrase:"}
		survey.AskOne(prompt, &passphrase)

		var newPassphrase string
		prompt = &survey.Password{Message: "Enter new passphrase:"}
		survey.AskOne(prompt, &newPassphrase)

		logger.Progress("Rotating passphrase for %s...\n", env)
		ctx := getContext()
		useCase := di.BuildRotatePassphrase()
		err = useCase.Execute(ctx, env, passphrase, newPassphrase)
		if err != nil {
			return err
		}

		clearCacheUseCase := di.BuildClearEnvCachedPassphrase()
		clearCacheUseCase.Execute(ctx, env)

		logger.Success("Passphrase rotated successfully")

		return nil
	},
}

func init() {
	rotateCmd.Flags().StringP("env", "e", "", "Environment Name")
	rotateCmd.MarkFlagRequired("env")

	rootCmd.AddCommand(rotateCmd)
}
