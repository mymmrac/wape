package io

import (
	"context"

	extism "github.com/extism/go-sdk"

	"github.com/mymmrac/wape/internal"
)

func Ready() extism.HostFunction {
	return internal.NewHostFunction("io.ready",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			handle := extism.DecodeI32(stack[0])

			result, ok := internal.IOHandles.GetOk(handle)
			if !ok {
				stack[0] = extism.EncodeI32(-1)
				return
			}

			stack[0] = extism.EncodeI32(result)

			if result != 0 {
				internal.IOHandles.Delete(handle)
			}
		}, []extism.ValueType{extism.ValueTypeI32}, []extism.ValueType{extism.ValueTypeI32},
	)
}
