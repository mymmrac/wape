package net

import (
	"context"
	"math/rand/v2"

	extism "github.com/extism/go-sdk"

	"github.com/mymmrac/wape/internal"
)

// ConnRead reads data from a connection.
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
		},
		[]extism.ValueType{extism.ValueTypeI32 /* connectionID */, extism.ValueTypePTR /* readDestination */},
		[]extism.ValueType{extism.ValueTypeI32 /* ioHandle | errorCode */},
	)
}

// ConnWrite writes data to a connection.
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
		},
		[]extism.ValueType{extism.ValueTypeI32 /* connectionID */, extism.ValueTypePTR /* writeSource */},
		[]extism.ValueType{extism.ValueTypeI32 /* ioHandle | errorCode */},
	)
}

// ConnClose closes the connection.
func ConnClose() extism.HostFunction {
	return internal.NewHostFunction("net.conn.close",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			connectionID := extism.DecodeI32(stack[0])

			conn := internal.Connections.Get(connectionID)

			var result int32
			if err := conn.Close(); err != nil {
				result = -1
			}

			stack[0] = extism.EncodeI32(result)
		},
		[]extism.ValueType{extism.ValueTypeI32 /* connectionID */},
		[]extism.ValueType{extism.ValueTypeI32 /* errorCode */},
	)
}
