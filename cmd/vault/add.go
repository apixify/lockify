package vault

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
	"github.com/spf13/cobra"
)

// AddCommand represents the add command for adding entries to the vault.
type AddCommand struct {
	useCase app.AddEntryUc
	prompt  service.PromptService
	logger  domain.Logger
	cmdCtx  *cli.CommandContext
}

// NewAddCommand creates a new add command instance.
func NewAddCommand(
	useCase app.AddEntryUc,
	prompt service.PromptService,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) (*cobra.Command, error) {
	cmd := &AddCommand{useCase, prompt, logger, cmdCtx}

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
	cobraCmd.Flags().BoolP(
		"secret",
		"s",
		false,
		"States that value to set is a secret and should be hidden in the terminal",
	)
	err := cobraCmd.MarkFlagRequired("env")
	if err != nil {
		return nil, fmt.Errorf("failed to mark env flag as required: %w", err)
	}

	return cobraCmd, nil
}

func (c *AddCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Progress("seting a new entry to the vault...")
	env, err := c.cmdCtx.RequireEnvFlag(cmd)
	if err != nil {
		return err
	}

	isSecret, err := cmd.Flags().GetBool("secret")
	if err != nil {
		c.logger.Error("failed to get secret flag: %w", err)
		return err
	}
	key, value, err := c.prompt.GetUserInputForKeyAndValue(isSecret)
	if err != nil {
		c.logger.Error("failed to get user input for key and value: %w", err)
		return err
	}

	shouldCache, err := c.cmdCtx.GetCacheFlag(cmd)
	if err != nil {
		c.logger.Error("failed to get cache flag: %w", err)
		return err
	}

	vctx := model.NewVaultContext(c.cmdCtx.GetContext(), env, shouldCache)
	dto := app.AddEntryDTO{Env: env, Key: key, Value: value}

	err = c.useCase.Execute(vctx, dto)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	c.logger.Success("key %s is added successfully.", key)

	return nil
}
