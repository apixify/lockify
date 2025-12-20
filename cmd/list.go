package cmd

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	useCase app.ListEntriesUc
	logger  domain.Logger
}

func NewListCommand(useCase app.ListEntriesUc, logger domain.Logger) *cobra.Command {
	cmd := &ListCommand{useCase, logger}
	// lockify list [env]
	cobraCmd := &cobra.Command{
		Use:   "list",
		Short: "List all keys in the vault",
		Long: `List all keys in the vault.

This command displays all keys stored in the vault for the specified environment.
Only keys are displayed, not decrypted values, for security reasons.`,
		Example: `  lockify list --env prod
  lockify list --env staging`,
		RunE: cmd.runE,
	}

	cobraCmd.Flags().StringP("env", "e", "", "Environment Name")
	cobraCmd.MarkFlagRequired("env")

	return cobraCmd
}

func (c *ListCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Progress("Listing all secrets in the vault")
	env, err := requireEnvFlag(cmd)
	if err != nil {
		return err
	}

	ctx := getContext()
	keys, err := c.useCase.Execute(ctx, env)
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		c.logger.Info("No entries found in vault")
		return nil
	}

	c.logger.Success("Found %d key(s):", len(keys))
	for _, v := range keys {
		c.logger.Output("  - %s\n", v)
	}

	return nil
}

func init() {
	listCmd := NewListCommand(di.BuildListEntries(), di.GetLogger())
	rootCmd.AddCommand(listCmd)
}
