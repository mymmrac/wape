package wasmgate

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	extism "github.com/extism/go-sdk"
	"github.com/tetratelabs/wazero"

	wio "github.com/mymmrac/wasm-gate/host/io"
	wnet "github.com/mymmrac/wasm-gate/host/net"
)

// Environment configures the behavior of WASM module.
//
// Note: Many environment configurations override each other, for example providing ModuleConfig will override all
// other configurations related to envs, args, FS, etc. Be aware that you may end up with unexpected behavior.
//
// TODO: Add json/yaml/toml tags
type Environment struct {
	// ==== Envs ====
	// Defaults to none.

	// Envs sets an environment variables.
	Envs []string
	// EnvsMap sets an environment variables as a map.
	EnvsMap map[string]string
	// EnvsFile sets an environment variables from a file.
	EnvsFile string
	// EnvsFromHost pass thought environment variables from the host.
	EnvsFromHost bool

	// ==== Args ====
	// Defaults to none.

	// Args assigns command-line arguments.
	Args []string
	// ArgsFile assigns command-line arguments from a file.
	ArgsFile string
	// ArgsFromHost pass thought command-line arguments from the host.
	ArgsFromHost bool

	// ==== Stdin ====
	// Defaults to always return io.EOF.

	// Stdin configures where standard input (file descriptor 0) is read.
	Stdin io.Reader
	// StdinFile configures standard input (file descriptor 0) to read from a file.
	StdinFile string
	// StdinFromHost pass thought stdin from the host.
	StdinFromHost bool

	// ==== Stdout ====
	// Defaults to io.Discard.

	// Stdout configures where standard output (file descriptor 1) is written.
	Stdout io.Writer
	// StdoutFile configures standard output (file descriptor 1) to write to a file.
	StdoutFile string
	// StdoutFromHost pass thought stdout from the host.
	StdoutFromHost bool

	// ==== Stderr ====
	// Defaults to io.Discard.

	// Stderr configures where standard error (file descriptor 2) is written.
	Stderr io.Writer
	// StderrFile configures standard error (file descriptor 2) to write to a file.
	StderrFile string
	// StderrFromHost pass thought stderr from the host.
	StderrFromHost bool

	// ==== FS ====
	// Defaults to have no file access.

	// FS configures the filesystem.
	FS fs.FS
	// FSFile configures the filesystem mount points.
	FSConfig wazero.FSConfig
	// FSDir configures the filesystem as a root directory.
	FSDir string
	// FSAllowedPaths configures the allowed filesystem paths that will be mapped in WASM module.
	FSAllowedPaths map[string]string
	// FSFromHost pass thought filesystem from the host.
	FSFromHost bool

	// TODO: Add file permissions + create/remove files

	// ==== Random Source ====
	// Defaults to a deterministic source.

	// RandSource configures a source of random bytes.
	RandSource io.Reader
	// RandSourceFile configures a source of random bytes from a file.
	RandSourceFile string
	// RandSourceFromHost pass thought random source from the host.
	RandSourceFromHost bool

	// ==== Wall Time ====
	// Defaults to always return UTC 1 January 1970.

	// WallTime returns the current unix/epoch time, seconds since midnight UTC 1 January 1970,
	// with a nanosecond fraction.
	WallTime func() (sec int64, ns int32)
	// WallTimeFromHost pass thought wall time from the host.
	WallTimeFromHost bool

	// ==== Nano Time ====
	// Defaults to always return 0.

	// NanoTime returns nanoseconds since an arbitrary start point, used to measure elapsed time.
	// This is sometimes referred to as a tick or monotonic time.
	NanoTime func() int64
	// NanoTimeFromHost pass thought nano time from the host.
	NanoTimeFromHost bool

	// ==== Nano Sleep ====
	// Defaults to always return immediately.

	// NanoSleep puts the current goroutine to sleep for at least ns nanoseconds.
	NanoSleep func(ns int64)
	// NanoSleepFromHost pass thought nano sleep from the host.
	NanoSleepFromHost bool

	// ==== Start Functions ====
	// Defaults to none.

	// StartFunctions configures the functions to call after the module is instantiated.
	StartFunctions []string

	// ==== Memory ====

	// MemoryLimitPages overrides the maximum pages allowed per memory. Defaults to 65536, allowing 4GB total memory
	// per instance if the maximum is not encoded in a Wasm binary. Max is 65536 (2^16) pages or 4GB.
	MemoryLimitPages uint32

	// MemoryCapacityFromMax eagerly allocates max memory. Defaults to false, which means minimum memory is
	// allocated and any call to grow memory results in re-allocations.
	MemoryCapacityFromMax bool

	// ==== Execution Time ====

	// MaxExecutionDuration limits the maximum execution time. Defaults to 0 (no limit).
	MaxExecutionDuration time.Duration

	// ==== Debug ====

	// DebugInfoEnabled toggles DWARF-based stack traces in the face of runtime errors. Defaults to false.
	DebugInfoEnabled bool

	// ExtismDebugEnvAllowed allows use of EXTISM_ENABLE_WASI_OUTPUT environment variable. Defaults to false (unsets
	// this variable, so it will not work for all WASM modules).
	ExtismDebugEnvAllowed bool

	// ==== Network ====

	// NetworkEnabled toggles network access. Defaults to false.
	NetworkEnabled bool

	// NetworksAllowed configures the allowed network protocols. See [net.Dial] for allowed protocols. Defaults to none.
	NetworksAllowed []string
	// NetworksAllowAll allows all network protocols. Defaults to false.
	NetworksAllowAll bool

	// NetworkAddressesAllowed configures the allowed network addresses. Defaults to none.
	NetworkAddressesAllowed []string
	// NetworkAddressesAllowAll allows all network addresses. Defaults to false.
	NetworkAddressesAllowAll bool

	// ==== WASI ====

	// DisableWASI disables WASI Preview 1 support. Defaults to false.
	DisableWASIP1 bool

	// ==== Wazero ===

	// ModuleConfig configures the WASM module.
	ModuleConfig wazero.ModuleConfig
	// RuntimeConfig configures the WASM runtime.
	RuntimeConfig wazero.RuntimeConfig

	// ==== Extism ====

	// Manifest is the plugin manifest that can be provided instead of configuration above.
	Manifest *extism.Manifest
	// PluginConfig is the plugin configuration that can be provided instead of configuration above.
	PluginConfig *extism.PluginConfig

	// ==== Host Functions ====

	// HostFunctions host functions available to the guest WASM module.
	HostFunctions []extism.HostFunction

	// ==== WASM Modules ====

	// Modules configures the WASM modules.
	Modules []ModuleData
}

// ModuleData represents WASM module with name and data.
type ModuleData struct {
	// Name of the WASM module.
	Name string

	// Source of WASM module.
	Data []byte

	// File path to read WASM module.
	File string

	// Url to download WASM module.
	Url string
	// Method to download WASM module.
	// Defaults to "GET".
	HttpMethod string
	// Headers to download WASM module.
	HttpHeaders map[string]string

	// SHA256 hash of WASM module to validate it.
	Hash string
}

// NewEnvironment returns a new environment.
func NewEnvironment() *Environment {
	return &Environment{}
}

// MakeModuleConfig returns the module configuration based on the environment.
func (e *Environment) MakeModuleConfig() wazero.ModuleConfig {
	if e.ModuleConfig != nil {
		return e.ModuleConfig
	}

	cfg := wazero.NewModuleConfig()

	if e.EnvsFromHost {
		for _, env := range os.Environ() {
			key, value, _ := strings.Cut(env, "=")
			cfg = cfg.WithEnv(key, value)
		}
	}

	if e.ArgsFromHost {
		cfg = cfg.WithArgs(os.Args...)
	}

	if e.StdinFromHost {
		cfg = cfg.WithStdin(os.Stdin)
	}

	if e.StdoutFromHost {
		cfg = cfg.WithStdout(os.Stdout)
	}

	if e.StderrFromHost {
		cfg = cfg.WithStderr(os.Stderr)
	}

	if e.FSFromHost {
		if runtime.GOOS == "windows" {
			wd, err := os.Getwd()
			if err != nil {
				cfg = cfg.WithFS(os.DirFS("C:\\"))
			} else {
				cfg = cfg.WithFS(os.DirFS(filepath.VolumeName(wd) + "\\"))
			}
		} else {
			cfg = cfg.WithFS(os.DirFS("/"))
		}
	}

	if e.RandSourceFromHost {
		cfg = cfg.WithRandSource(rand.Reader)
	}

	if e.WallTimeFromHost {
		cfg = cfg.WithSysWalltime()
	}

	if e.NanoTimeFromHost {
		cfg = cfg.WithSysNanotime()
	}

	if e.NanoSleepFromHost {
		cfg = cfg.WithSysNanosleep()
	}

	cfg.WithStartFunctions(e.StartFunctions...)

	return cfg
}

// MakeRuntimeConfig returns the runtime configuration based on the environment.
func (e *Environment) MakeRuntimeConfig() wazero.RuntimeConfig {
	if e.RuntimeConfig != nil {
		return e.RuntimeConfig
	}

	cfg := wazero.NewRuntimeConfig()

	cfg = cfg.WithDebugInfoEnabled(e.DebugInfoEnabled)

	if !e.ExtismDebugEnvAllowed {
		const env = "EXTISM_ENABLE_WASI_OUTPUT"
		err := os.Unsetenv(env)
		if err != nil {
			panic(fmt.Errorf("unset %q environment variable: %w", env, err))
		}
	}

	return cfg
}

// MakeManifest returns the manifest based on the environment.
func (e *Environment) MakeManifest() extism.Manifest {
	if e.Manifest != nil {
		return *e.Manifest
	}

	manifest := extism.Manifest{}

	for _, module := range e.Modules {
		var wasm extism.Wasm

		switch {
		case module.Data != nil:
			wasm = &extism.WasmData{
				Data: module.Data,
				Hash: module.Hash,
				Name: module.Name,
			}
		case module.File != "":
			wasm = &extism.WasmFile{
				Path: module.File,
				Hash: module.Hash,
				Name: module.Name,
			}
		case module.Url != "":
			wasm = &extism.WasmUrl{
				Url:     module.Url,
				Hash:    module.Hash,
				Headers: module.HttpHeaders,
				Name:    module.Name,
				Method:  module.HttpMethod,
			}
		default:
			continue
		}

		manifest.Wasm = append(manifest.Wasm, wasm)
	}

	return manifest
}

// MakePluginConfig returns the plugin configuration based on the environment.
func (e *Environment) MakePluginConfig() extism.PluginConfig {
	if e.PluginConfig != nil {
		return *e.PluginConfig
	}
	return extism.PluginConfig{
		RuntimeConfig: e.MakeRuntimeConfig(),
		EnableWasi:    !e.DisableWASIP1,
		ModuleConfig:  e.MakeModuleConfig(),
	}
}

// MakeHostFunctions returns the host functions based on the environment.
func (e *Environment) MakeHostFunctions() []extism.HostFunction {
	functions := make([]extism.HostFunction, 0, len(e.HostFunctions))

	if e.NetworkEnabled {
		functions = append(functions, wio.Ready())
	}

	if e.NetworkEnabled {
		functions = append(functions, wnet.Dial(wnet.DialConfig{
			NetworksAllowed:          e.NetworksAllowed,
			NetworksAllowAll:         e.NetworksAllowAll,
			NetworkAddressesAllowed:  e.NetworkAddressesAllowed,
			NetworkAddressesAllowAll: e.NetworkAddressesAllowAll,
		}))
		functions = append(functions, wnet.ConnRead())
		functions = append(functions, wnet.ConnWrite())
	}

	functions = append(functions, e.HostFunctions...)
	return functions
}
