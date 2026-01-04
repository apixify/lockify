package main

import (
	"log"

	"github.com/ahmed-abdelgawad92/lockify/cmd"
	"github.com/ahmed-abdelgawad92/lockify/cmd/cache"
	"github.com/ahmed-abdelgawad92/lockify/cmd/key"
	"github.com/ahmed-abdelgawad92/lockify/cmd/vault"
)

func main() {
	rootCmd, cmdCtx := cmd.NewRootCmd()

	if err := vault.RegisterCommands(rootCmd, cmdCtx); err != nil {
		log.Fatalf("Failed to register vault commands: %v", err)
	}

	if err := cache.RegisterCommands(rootCmd, cmdCtx); err != nil {
		log.Fatalf("Failed to register cache commands: %v", err)
	}

	if err := key.RegisterCommands(rootCmd, cmdCtx); err != nil {
		log.Fatalf("Failed to register key commands: %v", err)
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
