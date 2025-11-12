package cmd

import (
	"fmt"

	"github.com/apixify/lockify/internal/service"
	"github.com/spf13/cobra"
)

// lockify cache clear
var clearCmd = &cobra.Command{
	Use:   "cache clear",
	Short: "Clear cached passphrase.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("⏳ clearing cached passphrases")
		service.ClearAllPassphrases()
		fmt.Println("⏳ cleared cached passphrases")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
