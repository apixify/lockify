package key

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/service"
	"github.com/spf13/cobra"
)

// RotateCommand represents the key rotate command for rotating vault passphrases.
type RotateCommand struct {
	useCase app.RotatePassphraseUc
	prompt  service.PromptService
	logger  domain.Logger
	cmdCtx  *cli.CommandContext
}

// NewRotateCommand creates a new key rotate command instance.
func NewRotateCommand(
	useCase app.RotatePassphraseUc,
	prompt service.PromptService,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) (*cobra.Command, error) {
	cmd := &RotateCommand{useCase, prompt, logger, cmdCtx}

	// lockify key rotate --env [env]
	cobraCmd := &cobra.Command{
		Use:   "key rotate",
		Short: "Rotate the passphrase for a vault",
		Long: `Rotate the passphrase for a vault.

This command allows you to change the passphrase for a vault by re-encrypting all entries
with a new passphrase. You will be prompted for the current passphrase and a new passphrase.`,
		Example: `  lockify key rotate --env prod
  lockify key rotate --env staging`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment Name")
	err := cobraCmd.MarkFlagRequired("env")
	if err != nil {
		return nil, fmt.Errorf("failed to mark env flag as required: %w", err)
	}

	return cobraCmd, nil
}

func (c *RotateCommand) runE(cmd *cobra.Command, args []string) error {
	env, err := c.cmdCtx.RequireEnvFlag(cmd)
	if err != nil {
		return err
	}

	passphrase, err := c.prompt.GetPassphraseInput("Enter current passphrase:")
	if err != nil {
		return err
	}
	newPassphrase, err := c.prompt.GetPassphraseInput("Enter new passphrase:")
	if err != nil {
		return err
	}

	c.logger.Progress("Rotating passphrase for %s...\n", env)
	ctx := c.cmdCtx.GetContext()
	err = c.useCase.Execute(ctx, env, passphrase, newPassphrase)
	if err != nil {
		c.logger.Error("failed to rotate passphrase: %w", err)
		return err
	}

	clearCacheUseCase := di.BuildClearEnvCachedPassphrase()
	err = clearCacheUseCase.Execute(ctx, env)
	if err != nil {
		c.logger.Error("failed to clear cached passphrase: %w", err)
	}

	c.logger.Success("Passphrase rotated successfully")

	return nil
}
