package net

import (
	"context"
	"math/rand/v2"
	"net"

	extism "github.com/extism/go-sdk"

	wasmgate "github.com/mymmrac/wasm-gate"
	"github.com/mymmrac/wasm-gate/internal"
)

func Dial(_ *wasmgate.Environment) extism.HostFunction {
	return internal.NewHostFunction("net.dial",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			network, err := p.ReadString(stack[0])
			if err != nil {
				panic(err)
			}

			addr, err := p.ReadString(stack[1])
			if err != nil {
				panic(err)
			}

			// TODO: Validate network and address

			conn, err := net.Dial(network, addr)
			if err != nil {
				panic(err)
			}

			connID := rand.Int32N(10000)
			internal.Connections.Set(connID, conn)

			stack[0] = extism.EncodeI32(connID)
		}, []extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypeI32},
	)
}
