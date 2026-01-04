package cache

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/spf13/cobra"
)

// ClearCommand represents the clear command for clearing cached passphrases.
type ClearCommand struct {
	clearCachedPassphraseUc    app.ClearCachedPassphraseUc
	clearEnvCachedPassphraseUc app.ClearEnvCachedPassphraseUseCase
	logger                     domain.Logger
	cmdCtx                     *cli.CommandContext
}

// NewClearCommand creates a new clear command instance.
func NewClearCommand(
	clearCachedPassphraseUc app.ClearCachedPassphraseUc,
	clearEnvCachedPassphraseUc app.ClearEnvCachedPassphraseUseCase,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) *cobra.Command {
	cmd := &ClearCommand{
		clearCachedPassphraseUc,
		clearEnvCachedPassphraseUc,
		logger,
		cmdCtx,
	}

	// lockify cache clear [--env env]
	cobraCmd := &cobra.Command{
		Use:   "cache clear",
		Short: "Clear cached passphrases",
		Long: `Clear cached passphrases from the system keyring.

If --env flag is provided, only the passphrase for that environment is cleared.
Otherwise, all cached passphrases are cleared.`,
		Example: `  lockify cache clear
  lockify cache clear --env prod`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().
		StringP("env", "e", "", "Environment name (optional - clears specific env only)")

	return cobraCmd
}

func (c *ClearCommand) runE(cmd *cobra.Command, args []string) error {
	env, err := cmd.Flags().GetString("env")
	if err != nil {
		return fmt.Errorf("failed to get env flag: %w", err)
	}
	ctx := c.cmdCtx.GetContext()

	if env != "" {
		// Clear specific environment
		c.logger.Progress("Clearing cached passphrase for environment %q", env)
		err := c.clearEnvCachedPassphraseUc.Execute(model.NewVaultContext(ctx, env, false))
		if err != nil {
			c.logger.Error("failed to clear cached passphrase: %v", err)
			return err
		}
		c.logger.Success("Cleared cached passphrase for environment %q", env)
	} else {
		// Clear all
		c.logger.Progress("Clearing all cached passphrases")
		err := c.clearCachedPassphraseUc.Execute(model.NewVaultContext(ctx, env, false))
		if err != nil {
			c.logger.Error("failed to clear cached passphrases: %v", err)
			return err
		}
		c.logger.Success("Cleared all cached passphrases")
	}

	return nil
}
