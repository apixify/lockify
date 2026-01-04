package cmd

import (
	"github.com/ahmed-abdelgawad92/lockify/internal/domain"
	"github.com/spf13/cobra"
)

// Version holds the current version of the Lockify CLI.
var Version = "0.0.0"

// VersionCommand represents the version command for displaying the CLI version.
type VersionCommand struct {
	logger domain.Logger
}

// NewVersionCommand creates a new version command instance.
func NewVersionCommand(logger domain.Logger) *cobra.Command {
	cmd := &VersionCommand{logger}

	cobraCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the current version of Lockify",
		RunE:  cmd.runE,
	}

	return cobraCmd
}

func (c *VersionCommand) runE(cmd *cobra.Command, args []string) error {
	c.logger.Output("Lockify CLI v%s\n", Version)
	return nil
}
