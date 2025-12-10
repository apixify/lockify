package cmd

import (
	"fmt"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/spf13/cobra"
)

type ExportCommand struct {
	useCase app.ExportEnvUc
	logger  domain.Logger
}

func NewExportCommand(useCase app.ExportEnvUc, logger domain.Logger) *cobra.Command {
	cmd := &ExportCommand{useCase, logger}
	// lockify export --env [env] --format [dotenv|json]
	// lockify export --env prod --format dotenv > .env
	// lockify export --env staging --format json > env.json
	return &cobra.Command{
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
}

func (c *ExportCommand) runE(cmd *cobra.Command, args []string) error {
	env, err := requireEnvFlag(cmd)
	if err != nil {
		return err
	}

	format, err := requireStringFlag(cmd, "format")
	if err != nil {
		return err
	}

	expotFormat, err := value.NewFileFormat(format)
	if err != nil {
		return err
	}
	c.logger.Progress("Exporting entries for environment %s...\n", env)
	ctx := getContext()
	err = c.useCase.Execute(ctx, env, expotFormat)
	if err != nil {
		return fmt.Errorf("failed to export entries for environment %s: %w", env, err)
	}

	return nil
}

func init() {
	exportCmd := NewExportCommand(di.BuildExportEnv(), di.GetLogger())
	exportCmd.Flags().StringP("env", "e", "", "Environment Name")
	exportCmd.Flags().String("format", "dotenv", "The format of the exported file [dotenv|json]")
	exportCmd.MarkFlagRequired("env")

	rootCmd.AddCommand(exportCmd)
}
