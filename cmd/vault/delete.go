package vault

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/spf13/cobra"
)

// DeleteCommand represents the delete command for removing entries from the vault.
type DeleteCommand struct {
	useCase app.DeleteEntryUc
	logger  domain.Logger
	cmdCtx  *cli.CommandContext
}

// NewDeleteCommand creates a new delete command instance.
func NewDeleteCommand(
	useCase app.DeleteEntryUc,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) (*cobra.Command, error) {
	cmd := &DeleteCommand{useCase, logger, cmdCtx}
	// lockify del --env [env] --key [key]
	cobraCmd := &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "rm"},
		Short:   "Delete an entry from the vault",
		Long: `Delete an entry from the vault.

This command removes a key-value pair from the vault for the specified environment.`,
		Example: `  lockify delete --env prod --key OLD_KEY
  lockify del --env staging -k DEPRECATED_KEY`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment Name")
	cobraCmd.Flags().StringP("key", "k", "", "key to delete from the vault")
	err := cobraCmd.MarkFlagRequired("env")
	if err != nil {
		return nil, fmt.Errorf("failed to mark env flag as required: %w", err)
	}
	err = cobraCmd.MarkFlagRequired("key")
	if err != nil {
		return nil, fmt.Errorf("failed to mark key flag as required: %w", err)
	}

	return cobraCmd, nil
}

func (c *DeleteCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Progress("removing key...")
	env, err := c.cmdCtx.RequireEnvFlag(cmd)
	if err != nil {
		return err
	}
	key, err := c.cmdCtx.RequireStringFlag(cmd, "key")
	if err != nil {
		return err
	}

	ctx := c.cmdCtx.GetContext()
	err = c.useCase.Execute(ctx, env, key)
	if err != nil {
		return err
	}

	c.logger.Success("key %s is removed successfully.\n", key)

	return nil
}
