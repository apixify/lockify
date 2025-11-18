package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/apixify/lockify/internal/app"
	"github.com/apixify/lockify/internal/di"
	"github.com/spf13/cobra"
)

// lockify add --env [env]
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add or update an entry in the vault",
	Long: `Add or update an entry in the vault.

This command prompts you for a key and value, then encrypts and stores the value in the vault.
Use the --secret flag to hide the value input in the terminal.`,
	Example: `  lockify add --env prod
  lockify add --env staging --secret`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger.Progress("seting a new entry to the vault...")
		env, err := requireEnvFlag(cmd)
		if err != nil {
			return err
		}

		isSecret, _ := cmd.Flags().GetBool("secret")
		key, value := getUserInputForKeyAndValue(isSecret)

		ctx := getContext()
		useCase := di.BuildAddEntry()
		dto := app.AddEntryDTO{Env: env, Key: key, Value: value}

		err = useCase.Execute(ctx, dto)
		if err != nil {
			logger.Error(err.Error())
			return err
		}

		logger.Success("key %s is added successfully.", key)

		return nil
	},
}

func init() {
	addCmd.Flags().StringP("env", "e", "", "Environment Name")
	addCmd.Flags().BoolP("secret", "s", false, "States that value to set is a secret and should be hidden in the terminal")
	addCmd.MarkFlagRequired("env")

	rootCmd.AddCommand(addCmd)
}

func getUserInputForKeyAndValue(isSecret bool) (key, value string) {
	prompt := &survey.Input{Message: "Enter key:"}
	survey.AskOne(prompt, &key)

	if isSecret {
		prompt := &survey.Password{Message: "Enter secret:"}
		survey.AskOne(prompt, &value)
	} else {
		prompt = &survey.Input{Message: "Enter value:"}
		survey.AskOne(prompt, &value)
	}

	return key, value
}
