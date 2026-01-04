package vault

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/spf13/cobra"
)

// InitCommand represents the init command for initializing a new vault.
type InitCommand struct {
	useCase app.InitUc
	logger  domain.Logger
	cmdCtx  *cli.CommandContext
}

// NewInitCommand creates a new init command instance.
func NewInitCommand(
	initUc app.InitUc,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) (*cobra.Command, error) {
	cmd := &InitCommand{useCase: initUc, logger: logger, cmdCtx: cmdCtx}

	// lockify init --env [env]
	cobraCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Lockify vault in the current directory",
		Long: `Initialize a new Lockify vault for an environment.

	This command creates a new encrypted vault file that will store your environment variables.
	You will be prompted for a passphrase that will be used to encrypt and decrypt your secrets.`,
		Example: `  lockify init --env prod
	lockify init --env staging
	lockify init -e local`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment Name")
	err := cobraCmd.MarkFlagRequired("env")
	if err != nil {
		return nil, fmt.Errorf("failed to mark env flag as required: %w", err)
	}

	return cobraCmd, nil
}

func (c *InitCommand) runE(cobraCmd *cobra.Command, args []string) error {
	env, err := c.cmdCtx.RequireEnvFlag(cobraCmd)
	if err != nil {
		return err
	}

	shouldCache, err := c.cmdCtx.GetCacheFlag(cobraCmd)
	if err != nil {
		c.logger.Error("failed to get cache flag: %w", err)
		return err
	}

	c.logger.Progress("Initializing Lockify vault")
	vctx := model.NewVaultContext(c.cmdCtx.GetContext(), env, shouldCache)
	vault, err := c.useCase.Execute(vctx)
	if err != nil {
		return err
	}

	c.logger.Success("Lockify vault initialized at %s", vault.Path())
	return nil
}
