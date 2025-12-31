package vault

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/spf13/cobra"
)

// ListCommand represents the list command for listing vault entries.
type ListCommand struct {
	useCase app.ListEntriesUc
	logger  domain.Logger
	cmdCtx  *cli.CommandContext
}

// NewListCommand creates a new list command instance.
func NewListCommand(
	useCase app.ListEntriesUc,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) (*cobra.Command, error) {
	cmd := &ListCommand{useCase, logger, cmdCtx}
	// lockify list [env]
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys in the vault",
		Long: `List all keys in the vault.

This command displays all keys stored in the vault for the specified environment.
Only keys are displayed, not decrypted values, for security reasons.`,
		Example: `  lockify list --env prod
  lockify list --env staging`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment Name")
	err := cobraCmd.MarkFlagRequired("env")
	if err != nil {
		return nil, fmt.Errorf("failed to mark env flag as required: %w", err)
	}

	return cobraCmd, nil
}

func (c *ListCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Progress("Listing all secrets in the vault")
	env, err := c.cmdCtx.RequireEnvFlag(cmd)
	if err != nil {
		return err
	}

	ctx := c.cmdCtx.GetContext()
	keys, err := c.useCase.Execute(ctx, env)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		c.logger.Info("No entries found in vault")
		return nil
	}

	c.logger.Success("Found %d key(s):", len(keys))
	for _, v := range keys {
		c.logger.Output("  - %s\n", v)
	}

	return nil
}
