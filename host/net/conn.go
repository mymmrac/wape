package net

import (
	"context"
	"math/rand/v2"

	extism "github.com/extism/go-sdk"

	"github.com/mymmrac/wape/internal"
)

func ConnRead() extism.HostFunction {
	return internal.NewHostFunction("net.conn.read",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			connectionID := extism.DecodeI32(stack[0])

			length, err := p.Length(stack[1])
			if err != nil {
				panic(err)
			}

			buffer, ok := p.Memory().Read(uint32(stack[1]), uint32(length))
			if !ok {
				panic("failed to read buffer")
			}

			handle := rand.Int32()
			internal.IOHandles.Set(handle, 0)

			conn := internal.Connections.Get(connectionID)
			go func() {
				var n int
				n, err = conn.Read(buffer)
				if err != nil {
					panic(err)
				}
				internal.IOHandles.Set(handle, int32(n))
			}()

			stack[0] = extism.EncodeI32(handle)
		}, []extism.ValueType{extism.ValueTypeI32, extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypeI32},
	)
}

func ConnWrite() extism.HostFunction {
	return internal.NewHostFunction("net.conn.write",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			connectionID := extism.DecodeI32(stack[0])

			length, err := p.Length(stack[1])
			if err != nil {
				panic(err)
			}

			buffer, ok := p.Memory().Read(uint32(stack[1]), uint32(length))
			if !ok {
				panic("failed to read buffer")
			}

			handle := rand.Int32()
			internal.IOHandles.Set(handle, 0)

			conn := internal.Connections.Get(connectionID)
			go func() {
				var n int
				n, err = conn.Write(buffer)
				if err != nil {
					panic(err)
				}
				internal.IOHandles.Set(handle, int32(n))
			}()

			stack[0] = extism.EncodeI32(handle)
		}, []extism.ValueType{extism.ValueTypeI32, extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypeI32},
	)
}
