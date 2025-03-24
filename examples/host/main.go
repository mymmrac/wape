package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"

	wasmgate "github.com/mymmrac/wasm-gate"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	env := wasmgate.NewEnvironment()

	env.Modules = []wasmgate.ModuleData{
		{
			Name: "main",
			File: "./plugins/http/main.wasm",
		},
	}

	env.StdinFromHost = true
	env.StdoutFromHost = true

	env.FSFromHost = true

	env.WallTimeFromHost = true
	env.NanoTimeFromHost = true
	env.NanoSleepFromHost = true

	env.NetworkEnabled = true
	env.NetworksAllowAll = true
	env.NetworkAddressesAllowAll = true

	cmPlugin, err := wasmgate.NewCompiledPlugin(ctx, env)
	assert(err == nil, err)

	plugin, err := cmPlugin.Instance(ctx, env.MakePluginInstanceConfig())
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
