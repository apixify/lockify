package vault

import (
	"fmt"
	"os"

	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain/model/value"
	"github.com/spf13/cobra"
)

// ImportCommand represents the import command for importing entries into the vault.
type ImportCommand struct {
	useCase app.ImportEnvUc
	logger  domain.Logger
	cmdCtx  *cli.CommandContext
}

// NewImportCommand creates a new import command instance.
func NewImportCommand(
	useCase app.ImportEnvUc,
	logger domain.Logger,
	cmdCtx *cli.CommandContext,
) (*cobra.Command, error) {
	cmd := &ImportCommand{useCase, logger, cmdCtx}

	// lockify import .env --env prod --format dotenv
	// lockify import config.json --env staging --format json
	cobraCmd := &cobra.Command{
		Use:   "import [file]",
		Short: "Import variables from a file into the vault",
		Long: `Import variables from a file into the vault.

This command reads variables from a file and imports them into the vault.
Supported formats are dotenv (.env) and JSON.

If no file is specified, the command reads from stdin.`,
		Example: `  lockify import .env --env prod --format dotenv
  lockify import config.json --env staging --format json
  cat .env | lockify import --env local --format dotenv`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment name")
	cobraCmd.Flags().String("format", "dotenv", "Input format (dotenv|json)")
	cobraCmd.Flags().Bool("overwrite", false, "Overwrite existing keys")

	err := cobraCmd.MarkFlagRequired("env")
	if err != nil {
		return nil, fmt.Errorf("failed to mark env flag as required: %w", err)
	}
	err = cobraCmd.MarkFlagRequired("format")
	if err != nil {
		return nil, fmt.Errorf("failed to mark format flag as required: %w", err)
	}

	return cobraCmd, nil
}

func (c *ImportCommand) runE(cmd *cobra.Command, args []string) error {
	env, err := c.cmdCtx.RequireEnvFlag(cmd)
	if err != nil {
		return err
	}

	overwrite, err := cmd.Flags().GetBool("overwrite")
	if err != nil {
		c.logger.Error("failed to get overwrite flag: %w", err)
	}
	format, err := c.cmdCtx.RequireStringFlag(cmd, "format")
	if err != nil {
		return fmt.Errorf("failed to retrieve format flag: %w", err)
	}

	fileFormat, err := value.NewFileFormat(format)
	if err != nil {
		return err
	}

	file, filename, err := getFile(args)
	if err != nil {
		return err
	}

	if file != os.Stdin {
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				c.logger.Error("failed to close file: %v", closeErr)
			}
		}()
	}

	c.logger.Progress("Importing variables from %s...", filename)
	ctx := c.cmdCtx.GetContext()
	imported, skipped, err := c.useCase.Execute(ctx, env, fileFormat, file, overwrite)
	if err != nil {
		return fmt.Errorf("failed to import env variables: %w", err)
	}

	c.logger.Success("Imported %d key(s), skipped %d key(s)", imported, skipped)

	return nil
}

func getFile(args []string) (*os.File, string, error) {
	var file *os.File
	var filename string

	if len(args) > 0 {
		filename = args[0]
		var err error
		file, err = os.Open(filename)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open file %q: %w", filename, err)
		}
	} else {
		file = os.Stdin
		filename = "stdin"
	}

	return file, filename, nil
}
