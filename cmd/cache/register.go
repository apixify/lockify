package cache

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/cli"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/spf13/cobra"
)

// RegisterCommands registers all cache-related commands to the root command.
func RegisterCommands(rootCmd *cobra.Command, cmdCtx *cli.CommandContext) error {
	clearCmd := NewClearCommand(
		di.BuildClearCachedPassphrase(),
		di.BuildClearEnvCachedPassphrase(),
		di.GetLogger(),
		cmdCtx,
	)

	setCmd, err := NewSetCommand(
		di.BuildCachePassphrase(),
		di.BuildPromptService(),
		di.GetLogger(),
		cmdCtx,
	)
	if err != nil {
		return err
	}

	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(setCmd)

	return nil
}
