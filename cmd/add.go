package cmd

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
	"github.com/spf13/cobra"
)

type AddCommand struct {
	useCase app.AddEntryUc
	prompt  service.PromptService
	logger  domain.Logger
}

func NewAddCommand(useCase app.AddEntryUc, prompt service.PromptService, logger domain.Logger) *cobra.Command {
	cmd := &AddCommand{useCase, prompt, logger}

	// lockify add --env [env]
	cobraCmd := &cobra.Command{
		Use:   "add",
		Short: "Add or update an entry in the vault",
		Long: `Add or update an entry in the vault.

	This command prompts you for a key and value, then encrypts and stores the value in the vault.
	Use the --secret flag to hide the value input in the terminal.`,
		Example: `  lockify add --env prod
	lockify add --env staging --secret`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment Name")
	cobraCmd.Flags().BoolP("secret", "s", false, "States that value to set is a secret and should be hidden in the terminal")
	cobraCmd.MarkFlagRequired("env")

	return cobraCmd
}

func (c *AddCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Progress("seting a new entry to the vault...")
	env, err := requireEnvFlag(cmd)
	if err != nil {
		return err
	}

	isSecret, _ := cmd.Flags().GetBool("secret")
	key, value := c.prompt.GetUserInputForKeyAndValue(isSecret)

	ctx := getContext()
	dto := app.AddEntryDTO{Env: env, Key: key, Value: value}

	err = c.useCase.Execute(ctx, dto)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	c.logger.Success("key %s is added successfully.", key)

	return nil
}

func init() {
	addCmd := NewAddCommand(di.BuildAddEntry(), di.BuildPromptService(), di.GetLogger())
	rootCmd.AddCommand(addCmd)
}
