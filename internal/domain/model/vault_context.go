package model

import "context"

// VaultContext represents the context for a vault operation.
type VaultContext struct {
	Context     context.Context
	Env         string
	ShouldCache bool
}

// NewVaultContext creates a new VaultContext instance.
func NewVaultContext(ctx context.Context, env string, shouldCache bool) *VaultContext {
	return &VaultContext{Context: ctx, Env: env, ShouldCache: shouldCache}
}
