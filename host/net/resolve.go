package net

import (
	"context"
	"net"
	"strings"

	extism "github.com/extism/go-sdk"

	"github.com/mymmrac/wape/internal"
)

// LookupHost creates a host function that calls [net.Resolver.LookupHost].
func LookupHost() extism.HostFunction {
	return internal.NewHostFunction("net.resolver.lookupHost",
		func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			host, err := p.ReadString(stack[0])
			if err != nil {
				panic(err)
			}

			addresses, err := net.DefaultResolver.LookupHost(ctx, host)
			if err != nil {
				panic(err)
			}

			addressesPtr, err := p.WriteString(strings.Join(addresses, ","))
			if err != nil {
				panic(err)
			}

			stack[0] = addressesPtr
		},
		[]extism.ValueType{extism.ValueTypePTR /* host */},
		[]extism.ValueType{extism.ValueTypePTR /* addresses | errorCode */},
	)
}
