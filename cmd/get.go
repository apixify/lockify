package cmd

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/spf13/cobra"
)

type GetCommand struct {
	useCase app.GetEntryUc
	logger  domain.Logger
}

func NewGetCommand(useCase app.GetEntryUc, logger domain.Logger) *cobra.Command {
	cmd := &GetCommand{useCase, logger}
	// lockify get --env [env] --key [key]
	return &cobra.Command{
		Use:   "get",
		Short: "Get a decrypted value from the vault",
		Long: `Get a decrypted value from the vault.

This command retrieves and decrypts a value from the vault for the specified key.
The decrypted value is printed to stdout, making it suitable for shell scripting.`,
		Example: `  lockify get --env prod --key DATABASE_URL
  lockify get --env staging -k API_KEY`,
		RunE: cmd.runE,
	}
}

func (c *GetCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Progress("getting an entry from the vault")
	env, err := requireEnvFlag(cmd)
	if err != nil {
		return err
	}

	key, err := requireStringFlag(cmd, "key")
	if err != nil {
		return err
	}

	ctx := getContext()
	value, err := c.useCase.Execute(ctx, env, key)
	if err != nil {
		return err
	}

	c.logger.Success("retrieved key's value successfully")
	c.logger.Output(value)

	return nil
}

func init() {
	getCmd := NewGetCommand(di.BuildGetEntry(), di.GetLogger())
	getCmd.Flags().StringP("env", "e", "", "Environment name")
	getCmd.Flags().StringP("key", "k", "", "The key to use for getting the entry")
	getCmd.MarkFlagRequired("env")

	rootCmd.AddCommand(getCmd)
}
