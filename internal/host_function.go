package internal

import extism "github.com/extism/go-sdk"

func NewHostFunction(
	name string, callback extism.HostFunctionStackCallback, params []extism.ValueType, returnTypes []extism.ValueType,
) extism.HostFunction {
	f := extism.NewHostFunctionWithStack(name, callback, params, returnTypes)
	f.SetNamespace("wasm-gate:host/env")
	return f
}
