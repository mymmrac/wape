package wasmgate

import (
	"context"

	extism "github.com/extism/go-sdk"
)

// NewPlugin creates a new Extism plugin.
func NewPlugin(ctx context.Context, env *Environment) (*extism.Plugin, error) {
	return extism.NewPlugin(ctx, env.MakeManifest(), env.MakePluginConfig(), env.HostFunctions)
}

// NewCompiledPlugin creates a new compiled Extism plugin.
func NewCompiledPlugin(ctx context.Context, env *Environment) (*extism.CompiledPlugin, error) {
	return extism.NewCompiledPlugin(ctx, env.MakeManifest(), env.MakePluginConfig(), env.HostFunctions)
}
