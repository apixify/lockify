package vault

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/spf13/cobra"
)

// GetCommand represents the get command for retrieving entries from the vault.
type GetCommand struct {
	useCase app.GetEntryUc
	logger  domain.Logger
	cmdCtx  *cli.CommandContext
}

// NewGetCommand creates a new get command instance.
func NewGetCommand(
	useCase app.GetEntryUc,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) (*cobra.Command, error) {
	cmd := &GetCommand{useCase, logger, cmdCtx}
	// lockify get --env [env] --key [key]
	cobraCmd := &cobra.Command{
		Use:   "get",
		Short: "Get a decrypted value from the vault",
		Long: `Get a decrypted value from the vault.

This command retrieves and decrypts a value from the vault for the specified key.
The decrypted value is printed to stdout, making it suitable for shell scripting.`,
		Example: `  lockify get --env prod --key DATABASE_URL
  lockify get --env staging -k API_KEY`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment name")
	cobraCmd.Flags().StringP("key", "k", "", "The key to use for getting the entry")
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

func (c *GetCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Progress("getting an entry from the vault")
	env, err := c.cmdCtx.RequireEnvFlag(cmd)
	if err != nil {
		return err
	}

	key, err := c.cmdCtx.RequireStringFlag(cmd, "key")
	if err != nil {
		return err
	}

	ctx := c.cmdCtx.GetContext()
	value, err := c.useCase.Execute(ctx, env, key)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}

	c.logger.Success("retrieved key's value successfully")
	c.logger.Output(value)

	return nil
}
