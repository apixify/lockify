package cache

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
	"github.com/spf13/cobra"
)

// SetCommand represents the cache set command for caching passphrases.
type SetCommand struct {
	useCase app.CachePassphraseUc
	prompt  service.PromptService
	logger  domain.Logger
	cmdCtx  *cli.CommandContext
}

// NewSetCommand creates a new cache set command instance.
func NewSetCommand(
	useCase app.CachePassphraseUc,
	prompt service.PromptService,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) (*cobra.Command, error) {
	cmd := &SetCommand{useCase, prompt, logger, cmdCtx}

	// lockify cache set --env [env]
	cobraCmd := &cobra.Command{
		Use:   "cache set",
		Short: "Cache a passphrase for an environment",
		Long: `Cache a passphrase for an environment in the system keyring.

This command prompts for a passphrase, validates it against the vault,
and caches it in the system keyring for future use.`,
		Example: `  lockify cache set --env prod
  lockify cache set --env staging`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment name")
	err := cobraCmd.MarkFlagRequired("env")
	if err != nil {
		return nil, fmt.Errorf("failed to mark env flag as required: %w", err)
	}

	return cobraCmd, nil
}

func (c *SetCommand) runE(cmd *cobra.Command, args []string) error {
	env, err := c.cmdCtx.RequireEnvFlag(cmd)
	if err != nil {
		return err
	}

	c.logger.Progress("Prompting for passphrase...")

	// Get passphrase from user
	passphrase, err := c.prompt.GetPassphraseInput(
		fmt.Sprintf("Enter passphrase for environment %q:", env),
	)
	if err != nil {
		return fmt.Errorf("failed to get passphrase: %w", err)
	}

	if passphrase == "" {
		return fmt.Errorf("passphrase cannot be empty")
	}

	ctx := c.cmdCtx.GetContext()
	err = c.useCase.Execute(ctx, env, passphrase)
	if err != nil {
		c.logger.Error("failed to cache passphrase: %v", err)
		return err
	}

	c.logger.Success("Passphrase cached successfully for environment %q", env)
	return nil
}
