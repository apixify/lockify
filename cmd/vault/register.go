package vault

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/spf13/cobra"
)

// RegisterCommands registers all vault-related commands to the root command.
func RegisterCommands(rootCmd *cobra.Command, cmdCtx *cli.CommandContext) error {
	initCmd, err := NewInitCommand(di.BuildInitializeVault(), di.GetLogger(), cmdCtx)
	if err != nil {
		return err
	}

	addCmd, err := NewAddCommand(
		di.BuildAddEntry(),
		di.BuildPromptService(),
		di.GetLogger(),
		cmdCtx,
	)
	if err != nil {
		return err
	}

	getCmd, err := NewGetCommand(di.BuildGetEntry(), di.GetLogger(), cmdCtx)
	if err != nil {
		return err
	}

	listCmd, err := NewListCommand(di.BuildListEntries(), di.GetLogger(), cmdCtx)
	if err != nil {
		return err
	}

	deleteCmd, err := NewDeleteCommand(di.BuildDeleteEntry(), di.GetLogger(), cmdCtx)
	if err != nil {
		return err
	}

	exportCmd, err := NewExportCommand(di.BuildExportEnv(), di.GetLogger(), cmdCtx)
	if err != nil {
		return err
	}

	importCmd, err := NewImportCommand(di.BuildImportEnv(), di.GetLogger(), cmdCtx)
	if err != nil {
		return err
	}

	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(importCmd)

	return nil
}
