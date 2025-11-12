package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// lockify list [env]
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Displays only keys, not decrypted values.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Listing all secrets in the vault")
		env, _ := cmd.Flags().GetString("env")
		vaultPath := filepath.Join(".lockify", env+".vault.enc")
		fmt.Println("Environment:" + env)
		fmt.Println("Vault path:" + vaultPath)
		return nil
	},
}

func init() {
	listCmd.Flags().String("env", defaultEnv, "The environment for which to list the secrets")

	rootCmd.AddCommand(listCmd)
}
