package cmd

import (
	"fmt"
	"os"

	"github.com/apixify/lockify/internal/di"
	"github.com/apixify/lockify/internal/domain/model/value"
	"github.com/spf13/cobra"
)

// lockify import .env --env prod --format dotenv
// lockify import config.json --env staging --format json
var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import variables from a file into the vault",
	Long: `Import variables from a file into the vault.

This command reads variables from a file and imports them into the vault.
Supported formats are dotenv (.env) and JSON.

If no file is specified, the command reads from stdin.`,
	Example: `  lockify import .env --env prod --format dotenv
  lockify import config.json --env staging --format json
  cat .env | lockify import --env local --format dotenv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		env, err := requireEnvFlag(cmd)
		if err != nil {
			return err
		}

		overwrite, _ := cmd.Flags().GetBool("overwrite")
		format, err := requireStringFlag(cmd, "format")
		if err != nil {
			return fmt.Errorf("failed to retrieve format flag: %w", err)
		}

		fileFormat := value.NewFileFormat(format)
		if !fileFormat.IsValid() {
			return fmt.Errorf("format must be either %q or %q, got %q", value.Json, value.DotEnv, format)
		}

		file, filename, err := getFile(args)
		if err != nil {
			return err
		}

		logger.Progress("Importing variables from %s...", filename)
		ctx := getContext()
		useCase := di.BuildImportEnv()
		imported, skipped, err := useCase.Execute(ctx, env, fileFormat, file, overwrite)
		if err != nil {
			return fmt.Errorf("failed to import env variables: %w", err)
		}

		logger.Success("Imported %d key(s), skipped %d key(s)", imported, skipped)

		return nil
	},
}

func init() {
	importCmd.Flags().StringP("env", "e", "", "Environment name")
	importCmd.Flags().String("format", "dotenv", "Input format (dotenv|json)")
	importCmd.Flags().Bool("overwrite", false, "Overwrite existing keys")

	importCmd.MarkFlagFilename("file")
	importCmd.MarkFlagRequired("env")
	importCmd.MarkFlagRequired("format")

	rootCmd.AddCommand(importCmd)
}

func getFile(args []string) (*os.File, string, error) {
	var file *os.File
	var filename string

	if len(args) > 0 {
		filename = args[0]
		file, err := os.Open(filename)
		if err != nil {
			return nil, "", fmt.Errorf("failed to open file %q: %w", filename, err)
		}
		defer file.Close()
	} else {
		file = os.Stdin
		filename = "stdin"
	}

	return file, filename, nil
}
