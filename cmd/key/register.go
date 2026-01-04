package key

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/spf13/cobra"
)

// RegisterCommands registers all key-related commands to the root command.
func RegisterCommands(rootCmd *cobra.Command, cmdCtx *cli.CommandContext) error {
	rotateCmd, err := NewRotateCommand(
		di.BuildRotatePassphrase(),
		di.BuildPromptService(),
		di.GetLogger(),
		cmdCtx,
	)
	if err != nil {
		return err
	}

	rootCmd.AddCommand(rotateCmd)

	return nil
}
