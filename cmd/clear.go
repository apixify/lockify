package cmd

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/spf13/cobra"
)

type ClearCommand struct {
	buildClearCachedPassphraseUc func() app.ClearCachedPassphraseUc
	logger                       domain.Logger
}

func NewClearCommand(buildClearCachedPassphraseUc func() app.ClearCachedPassphraseUc, logger domain.Logger) *cobra.Command {
	cmd := &ClearCommand{buildClearCachedPassphraseUc, logger}

	// lockify cache clear
	return &cobra.Command{
		Use:   "cache clear",
		Short: "Clear cached passphrases",
		Long: `Clear all cached passphrases from the system keyring.

This command removes all passphrases that were cached in the system keyring.
You will be prompted for passphrases again on next use.`,
		Example: `  lockify cache clear`,
		RunE:    cmd.runE,
	}
}

func (c *ClearCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Progress("clearing cached passphrases")
	useCase := c.buildClearCachedPassphraseUc()

	ctx := getContext()
	err := useCase.Execute(ctx)
	if err != nil {
		c.logger.Error("failed to cleare cached passphrases")
		return err
	}

	c.logger.Success("cleared cached passphrases")
	return nil
}

func init() {
	clearCmd := NewClearCommand(di.BuildClearCachedPassphrase, di.GetLogger())
	rootCmd.AddCommand(clearCmd)
}
