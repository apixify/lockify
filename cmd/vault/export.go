package vault

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/spf13/cobra"
)

// ExportCommand represents the export command for exporting vault entries.
type ExportCommand struct {
	useCase app.ExportEnvUc
	logger  domain.Logger
	cmdCtx  *cli.CommandContext
}

// NewExportCommand creates a new export command instance.
func NewExportCommand(
	useCase app.ExportEnvUc,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) (*cobra.Command, error) {
	cmd := &ExportCommand{useCase, logger, cmdCtx}
	// lockify export --env [env] --format [dotenv|json]
	// lockify export --env prod --format dotenv > .env
	// lockify export --env staging --format json > env.json
	cobraCmd := &cobra.Command{
		Use:   "export",
		Short: "Export all decrypted variables in a specific format",
		Long: `Export all decrypted variables in a specific format.

This command decrypts all entries in the vault and exports them in the specified format.
Use stdout redirection to save to a file (e.g., lockify export --env prod --format dotenv > .env).

The output is written to stdout, making it suitable for shell redirection.`,
		Example: `  lockify export --env prod --format dotenv > .env
  lockify export --env staging --format json > env.json
  lockify export --env local --format dotenv`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment Name")
	cobraCmd.Flags().String("format", "dotenv", "The format of the exported file [dotenv|json]")
	err := cobraCmd.MarkFlagRequired("env")
	if err != nil {
		return nil, fmt.Errorf("failed to mark env flag as required: %w", err)
	}

	return cobraCmd, nil
}

func (c *ExportCommand) runE(cmd *cobra.Command, args []string) error {
	env, err := c.cmdCtx.RequireEnvFlag(cmd)
	if err != nil {
		return err
	}

	format, err := c.cmdCtx.RequireStringFlag(cmd, "format")
	if err != nil {
		return err
	}

	expotFormat, err := value.NewFileFormat(format)
	if err != nil {
		return err
	}
	c.logger.Progress("Exporting entries for environment %s...\n", env)
	shouldCache, err := c.cmdCtx.GetCacheFlag(cmd)
	if err != nil {
		c.logger.Error("failed to get cache flag: %w", err)
		return err
	}

	vctx := model.NewVaultContext(c.cmdCtx.GetContext(), env, shouldCache)
	err = c.useCase.Execute(vctx, expotFormat)
	if err != nil {
		return fmt.Errorf("failed to export entries for environment %s: %w", env, err)
	}

	return nil
}
