package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/mymmrac/wape"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	env := wape.NewEnvironment()

	env.Modules = []wape.ModuleData{
		{
			Name: "main",
			File: "./plugins/http/main.wasm",
		},
	}

	env.CompilationCacheDir = "/tmp/wazero-cache"

	env.StdinFromHost = true
	env.StdoutFromHost = true

	env.FSFromHost = true

	env.WallTimeFromHost = true
	env.NanoTimeFromHost = true
	env.NanoSleepFromHost = true

	env.NetworkEnabled = true
	env.NetworksAllowAll = true
	env.NetworkAddressesAllowAll = true

	start := time.Now()
	cmPlugin, err := wape.NewCompiledPlugin(ctx, env)
	fmt.Println("Time:", time.Since(start))

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
