package cmd

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/spf13/cobra"
)

type DeleteCommand struct {
	useCase app.DeleteEntryUc
	logger  domain.Logger
}

func NewDeleteCommand(useCase app.DeleteEntryUc, logger domain.Logger) *cobra.Command {
	cmd := &DeleteCommand{useCase, logger}
	// lockify del --env [env] --key [key]
	return &cobra.Command{
		Use:     "delete",
		Aliases: []string{"del", "rm"},
		Short:   "Delete an entry from the vault",
		Long: `Delete an entry from the vault.

This command removes a key-value pair from the vault for the specified environment.`,
		Example: `  lockify delete --env prod --key OLD_KEY
  lockify del --env staging -k DEPRECATED_KEY`,
		RunE: cmd.runE,
	}
}

func (c *DeleteCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Progress("removing key...")
	env, err := requireEnvFlag(cmd)
	if err != nil {
		return err
	}
	key, err := requireStringFlag(cmd, "key")
	if err != nil {
		return err
	}

	ctx := getContext()
	err = c.useCase.Execute(ctx, env, key)
	if err != nil {
		return err
	}

	c.logger.Success("key %s is removed successfully.\n", key)

	return nil
}

func init() {
	delCmd := NewDeleteCommand(di.BuildDeleteEntry(), di.GetLogger())
	delCmd.Flags().StringP("env", "e", "", "Environment Name")
	delCmd.Flags().StringP("key", "k", "", "key to delete from the vault")
	delCmd.MarkFlagRequired("env")
	delCmd.MarkFlagRequired("key")

	rootCmd.AddCommand(delCmd)
}
