package cmd

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/app"
	"github.com/ahmed-abdelgawad92/lockify/internal/di"
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/spf13/cobra"
)

type InitCommand struct {
	useCase app.InitUc
	logger  domain.Logger
}

func NewInitCommand(initUc app.InitUc, logger domain.Logger) *cobra.Command {
	cmd := &InitCommand{useCase: initUc, logger: logger}

	// lockify init --env [env]
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Lockify vault in the current directory",
		Long: `Initialize a new Lockify vault for an environment.

	This command creates a new encrypted vault file that will store your environment variables.
	You will be prompted for a passphrase that will be used to encrypt and decrypt your secrets.`,
		Example: `  lockify init --env prod
	lockify init --env staging
	lockify init -e local`,
		RunE: cmd.runE,
	}
}

func (c *InitCommand) runE(cmd *cobra.Command, args []string) error {
	env, err := requireEnvFlag(cmd)
	if err != nil {
		return err
	}

	c.logger.Progress("Initializing Lockify vault")
	ctx := getContext()
	vault, err := c.useCase.Execute(ctx, env)
	if err != nil {
		return err
	}

	c.logger.Success("Lockify vault initialized at %s", vault.Path())
	return nil
}

func init() {
	initCmd := NewInitCommand(di.BuildInitializeVault(), di.GetLogger())
	initCmd.Flags().StringP("env", "e", "", "Environment Name")
	initCmd.MarkFlagRequired("env")

	rootCmd.AddCommand(initCmd)
}
