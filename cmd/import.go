package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// lockify import .env --env prod --format dotenv
// lockify import config.json --env staging --format json
var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "import a secret from the vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("importing a secret from the vault")
		return nil
	},
}

func init() {
	importCmd.Flags().StringP("env", "e", "", "Environment name")
	importCmd.Flags().StringP("format", "f", "", "Input format (dotenv|json)")
	importCmd.Flags().Bool("overwrite", false, "Overwrite existing keys")

	rootCmd.AddCommand(importCmd)
}
