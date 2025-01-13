package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	extism "github.com/extism/go-sdk"

	wasmgate "github.com/mymmrac/wasm-gate"
	"github.com/mymmrac/wasm-gate/host/io"
	"github.com/mymmrac/wasm-gate/host/net"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	env := wasmgate.NewEnvironment()
	env.Manifest = &extism.Manifest{
		Wasm: []extism.Wasm{
			&extism.WasmFile{
				Path: "./plugins/http/main.wasm",
			},
		},
	}

	env.StdinFromHost = true
	env.StdoutFromHost = true
	env.FSFromHost = true
	env.WallTimeFromHost = true
	env.NanoTimeFromHost = true
	env.NanoSleepFromHost = true

	env.HostFunctions = []extism.HostFunction{
		io.Ready(env),
		net.Dial(env),
		net.ConnRead(env),
		net.ConnWrite(env),
	}

	plugin, err := wasmgate.NewPlugin(ctx, env)
	assert(err == nil, err)

	exit, _, err := plugin.CallWithContext(ctx, "main", nil)
	assert(err == nil, err)
	fmt.Println("Exit code:", exit)

	err = plugin.Close(ctx)
	assert(err == nil, err)
}

func assert(ok bool, args ...any) {
	if !ok {
		_, file, line, _ := runtime.Caller(1)
		panic(fmt.Errorf("asser: %s:%d: %s", file, line, fmt.Sprint(args...)))
	}
}
