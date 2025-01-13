package main

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sync"

	extism "github.com/extism/go-sdk"

	wasmgate "github.com/mymmrac/wasm-gate"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	env := wasmgate.NewEnvironment()
	env.Manifest = &extism.Manifest{
		Wasm: []extism.Wasm{
			&extism.WasmFile{
				Path: "./examples/plugins/http/main.wasm",
			},
		},
	}

	env.StdinFromHost = true
	env.StdoutFromHost = true
	env.FSFromHost = true
	env.WallTimeFromHost = true
	env.NanoTimeFromHost = true
	env.NanoSleepFromHost = true

	ioHandles := NewSyncMap[int32, int32]()
	connections := NewSyncMap[int32, net.Conn]()

	env.HostFunctions = []extism.HostFunction{
		extism.NewHostFunctionWithStack("dial", func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
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

			connectionID := rand.Int32N(10000)
			connections.Set(connectionID, conn)

			stack[0] = extism.EncodeI32(connectionID)
		}, []extism.ValueType{extism.ValueTypePTR, extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypeI32}),
		extism.NewHostFunctionWithStack("read", func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			connectionID := extism.DecodeI32(stack[0])

			length, err := p.Length(stack[1])
			if err != nil {
				panic(err)
			}

			buffer, ok := p.Memory().Read(uint32(stack[1]), uint32(length))
			if !ok {
				panic("failed to read buffer")
			}

			ioHandle := rand.Int32()
			ioHandles.Set(ioHandle, 0)

			conn := connections.Get(connectionID)
			go func() {
				var n int
				n, err = conn.Read(buffer)
				if err != nil {
					panic(err)
				}
				ioHandles.Set(ioHandle, int32(n))
			}()

			stack[0] = extism.EncodeI32(ioHandle)
		}, []extism.ValueType{extism.ValueTypeI32, extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypeI32}),
		extism.NewHostFunctionWithStack("write", func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			connectionID := extism.DecodeI32(stack[0])

			length, err := p.Length(stack[1])
			if err != nil {
				panic(err)
			}

			buffer, ok := p.Memory().Read(uint32(stack[1]), uint32(length))
			if !ok {
				panic("failed to read buffer")
			}

			ioHandle := rand.Int32()
			ioHandles.Set(ioHandle, 0)

			conn := connections.Get(connectionID)
			go func() {
				var n int
				n, err = conn.Write(buffer)
				if err != nil {
					panic(err)
				}
				ioHandles.Set(ioHandle, int32(n))
			}()

			stack[0] = extism.EncodeI32(ioHandle)
		}, []extism.ValueType{extism.ValueTypeI32, extism.ValueTypePTR}, []extism.ValueType{extism.ValueTypeI32}),
		extism.NewHostFunctionWithStack("io_ready", func(ctx context.Context, p *extism.CurrentPlugin, stack []uint64) {
			ioHandle := extism.DecodeI32(stack[0])

			result, ok := ioHandles.GetOk(ioHandle)
			if !ok {
				stack[0] = extism.EncodeI32(-1)
				return
			}

			stack[0] = extism.EncodeI32(result)

			if result != 0 {
				ioHandles.Delete(ioHandle)
			}
		}, []extism.ValueType{extism.ValueTypeI32}, []extism.ValueType{extism.ValueTypeI32}),
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

type SyncMap[K comparable, V any] struct {
	m map[K]V
	l sync.RWMutex
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: make(map[K]V),
		l: sync.RWMutex{},
	}
}

func (m *SyncMap[K, V]) Get(key K) V {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m[key]
}

func (m *SyncMap[K, V]) GetOk(key K) (V, bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	value, ok := m.m[key]
	return value, ok
}

func (m *SyncMap[K, V]) Set(key K, value V) {
	m.l.Lock()
	defer m.l.Unlock()
	m.m[key] = value
}

func (m *SyncMap[K, _]) Delete(key K) {
	m.l.Lock()
	delete(m.m, key)
	m.l.Unlock()
}
